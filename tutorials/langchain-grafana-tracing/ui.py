"""Web UI for LangChain tracing and trace comparison.

Provides a Flask web server with HTML forms for testing question answering
and trace comparison workflows.
"""

from __future__ import annotations

import atexit
import json
import os
import re
from urllib.error import HTTPError, URLError
from urllib.request import Request, urlopen

from flask import Flask, render_template, request, jsonify
from langchain_core.messages import AIMessage, ToolMessage

from graph import AgentGraph
from incident_agent import IncidentAgent
from loki import emit_log, fetch_logs_for_trace
from config import Config
from tempo import fetch_trace_from_tempo
from trace_analyzer import compare_traces
from tracing import configure_tracing, get_tracer

app = Flask(__name__, template_folder="templates")
TRACER = None
OTEL_PROVIDER = None
AGENT_GRAPH = None
INCIDENT_AGENT = None

TRACE_ID_PATTERN = re.compile(r"^[0-9a-fA-F]{32}$")
MAX_QUESTION_LENGTH = 4000


def _json_body():
    """Return a JSON object or a useful client error response."""
    data = request.get_json(silent=True)
    if not isinstance(data, dict):
        return None, (jsonify({
            "success": False,
            "code": "invalid_json",
            "error": "Request body must be a JSON object",
        }), 400)
    return data, None


def _required_text(data, key: str, *, max_length: int = MAX_QUESTION_LENGTH):
    value = data.get(key)
    if not isinstance(value, str) or not value.strip():
        return None, f"{key} is required"
    value = value.strip()
    if len(value) > max_length:
        return None, f"{key} must be {max_length} characters or fewer"
    return value, None


def _service_status(name: str, url: str) -> dict:
    """Perform a short, read-only readiness check for a dependency."""
    try:
        with urlopen(Request(url, headers={"Accept": "application/json"}), timeout=1.5) as response:
            return {"name": name, "available": 200 <= response.status < 500}
    except (HTTPError, URLError, TimeoutError, OSError):
        return {"name": name, "available": False}


def _incident_payload(result: dict) -> dict:
    """Turn LangChain messages into bounded, UI-friendly investigation evidence."""
    tools = []
    evidence = []
    answer = "The investigation completed without a final response."
    for message in result.get("messages", []):
        if isinstance(message, AIMessage):
            for call in message.tool_calls:
                name = call.get("name", "unknown")
                if name not in tools:
                    tools.append(name)
            if message.content:
                answer = message.content if isinstance(message.content, str) else str(message.content)
        elif isinstance(message, ToolMessage):
            if message.name and message.name not in tools:
                tools.append(message.name)
            try:
                content = json.loads(message.content)
            except (TypeError, json.JSONDecodeError):
                content = {"ok": False, "message": str(message.content)[:1000]}
            evidence.append({"tool": message.name or "tool", "data": content})
    return {"answer": answer, "selected_tools": tools, "evidence": evidence}


@app.after_request
def disable_ui_caching(response):
    """Ensure rebuilt UI templates are not hidden by browser caches."""
    if response.content_type and response.content_type.startswith("text/html"):
        response.headers["Cache-Control"] = "no-store, max-age=0"
        response.headers["Pragma"] = "no-cache"
        response.headers["Expires"] = "0"
    return response


@app.route("/", methods=["GET"])
def index():
    """Main UI page."""
    return render_template("index.html")


@app.route("/api/question", methods=["POST"])
def ask_question():
    """Handle question answering endpoint."""
    global TRACER, AGENT_GRAPH
    if TRACER is None:
        configure_tracing()
        TRACER = get_tracer()

    data, error_response = _json_body()
    if error_response:
        return error_response
    question, validation_error = _required_text(data, "question")
    if validation_error:
        return jsonify({"success": False, "code": "invalid_question", "error": validation_error}), 400
    thread_id = str(data.get("thread_id", "default")).strip()[:128] or "default"

    try:
        with TRACER.start_as_current_span("langchain.request") as span:
            span.set_attribute("app.request.type", "web_question")
            span.set_attribute("app.thread.id", thread_id)
            trace_id = span.get_span_context().trace_id
            emit_log("request.started", thread_id=thread_id)
            result = AGENT_GRAPH.invoke(question, thread_id)
            emit_log("request.completed", thread_id=thread_id)

        if OTEL_PROVIDER is not None:
            OTEL_PROVIDER.force_flush()

        return jsonify({
            "success": True,
            "answer": result["answer"],
            "trace_id": f"{trace_id:032x}",
            "thread_id": thread_id,
            "sources": list(dict.fromkeys(
                document["source"] for document in result.get("documents", [])
            )),
        })
    except Exception as e:
        return jsonify({
            "success": False,
            "code": "question_failed",
            "error": str(e),
        }), 500


@app.route("/api/investigate", methods=["POST"])
def investigate_incident():
    """Run the bounded read-only incident agent and expose its evidence."""
    global TRACER, INCIDENT_AGENT
    if TRACER is None:
        configure_tracing()
        TRACER = get_tracer()
    data, error_response = _json_body()
    if error_response:
        return error_response
    question, validation_error = _required_text(data, "question")
    if validation_error:
        return jsonify({"success": False, "code": "invalid_question", "error": validation_error}), 400
    thread_id = str(data.get("thread_id", "incident-web")).strip()[:128] or "incident-web"
    try:
        max_model_calls = max(1, min(int(data.get("max_model_calls", 6)), 12))
    except (TypeError, ValueError):
        return jsonify({
            "success": False,
            "code": "invalid_model_calls",
            "error": "max_model_calls must be an integer between 1 and 12",
        }), 400

    try:
        with TRACER.start_as_current_span("incident.investigation") as span:
            span.set_attribute("app.request.type", "web_investigation")
            span.set_attribute("app.thread.id", thread_id)
            span.set_attribute("agent.max_model_calls", max_model_calls)
            trace_id = span.get_span_context().trace_id
            result = INCIDENT_AGENT.invoke(question, thread_id, max_model_calls)
            payload = _incident_payload(result)
        if OTEL_PROVIDER is not None:
            OTEL_PROVIDER.force_flush()
        return jsonify({
            "success": True,
            "trace_id": f"{trace_id:032x}",
            "thread_id": thread_id,
            **payload,
        })
    except Exception as error:
        return jsonify({
            "success": False,
            "code": "investigation_failed",
            "error": str(error),
        }), 500


@app.route("/api/compare", methods=["POST"])
def compare_traces_api():
    """Handle trace comparison endpoint."""
    data, error_response = _json_body()
    if error_response:
        return error_response
    trace_a_id = str(data.get("trace_a", "")).strip().lower()
    trace_b_id = str(data.get("trace_b", "")).strip().lower()

    if not TRACE_ID_PATTERN.fullmatch(trace_a_id) or not TRACE_ID_PATTERN.fullmatch(trace_b_id):
        return jsonify({
            "success": False,
            "code": "invalid_trace_ids",
            "error": "Both trace IDs must contain exactly 32 hexadecimal characters",
        }), 400

    try:
        trace_a = fetch_trace_from_tempo(trace_a_id)
        trace_b = fetch_trace_from_tempo(trace_b_id)
        result = compare_traces(trace_a, trace_b)
        comparison = result.to_dict()
        comparison["trace_a_logs"] = fetch_logs_for_trace(trace_a_id)
        comparison["trace_b_logs"] = fetch_logs_for_trace(trace_b_id)

        return jsonify({
            "success": True,
            "comparison": comparison,
        })
    except Exception as e:
        return jsonify({
            "success": False,
            "code": "comparison_failed",
            "error": str(e),
        }), 500


@app.route("/api/status", methods=["GET"])
def service_status():
    """Report dependency reachability for the UI connection indicators."""
    checks = [
        _service_status("Tempo", f"{Config.TEMPO_ENDPOINT}/ready"),
        _service_status("Loki", f"{Config.LOKI_ENDPOINT}/ready"),
        _service_status("Ollama", f"{os.getenv('OLLAMA_BASE_URL', 'http://localhost:11434').rstrip('/')}/api/tags"),
    ]
    return jsonify({"success": True, "services": checks})


@app.route("/health", methods=["GET"])
def health():
    """Health check endpoint."""
    return jsonify({"status": "healthy"}), 200


def create_app():
    """Create and configure the Flask application."""
    provider = configure_tracing()
    global TRACER, OTEL_PROVIDER, AGENT_GRAPH, INCIDENT_AGENT
    TRACER = get_tracer()
    OTEL_PROVIDER = provider
    AGENT_GRAPH = AgentGraph()
    INCIDENT_AGENT = IncidentAgent()

    def shutdown_tracing():
        """Flush and shut down tracing when the process exits."""
        provider.force_flush()
        provider.shutdown()
        AGENT_GRAPH.close()
        INCIDENT_AGENT.close()

    atexit.register(shutdown_tracing)

    return app


if __name__ == "__main__":
    app = create_app()
    app.run(host="0.0.0.0", port=5000, debug=False)
