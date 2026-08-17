"""Web UI for LangChain tracing and trace comparison.

Provides a Flask web server with HTML forms for testing question answering
and trace comparison workflows.
"""

from __future__ import annotations

import atexit

from flask import Flask, render_template, request, jsonify

from graph import AgentGraph
from loki import emit_log, fetch_logs_for_trace
from tempo import fetch_trace_from_tempo
from trace_analyzer import compare_traces
from tracing import configure_tracing, get_tracer

app = Flask(__name__, template_folder="templates")
TRACER = None
OTEL_PROVIDER = None
AGENT_GRAPH = None


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

    data = request.get_json()
    question = data.get("question", "What problem does LangChain solve?")
    thread_id = str(data.get("thread_id", "default"))[:128]

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
            "error": str(e),
        }), 500


@app.route("/api/compare", methods=["POST"])
def compare_traces_api():
    """Handle trace comparison endpoint."""
    data = request.get_json()
    trace_a_id = data.get("trace_a")
    trace_b_id = data.get("trace_b")

    if not trace_a_id or not trace_b_id:
        return jsonify({
            "success": False,
            "error": "Both trace_a and trace_b are required",
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
            "error": str(e),
        }), 500


@app.route("/health", methods=["GET"])
def health():
    """Health check endpoint."""
    return jsonify({"status": "healthy"}), 200


def create_app():
    """Create and configure the Flask application."""
    provider = configure_tracing()
    global TRACER, OTEL_PROVIDER, AGENT_GRAPH
    TRACER = get_tracer()
    OTEL_PROVIDER = provider
    AGENT_GRAPH = AgentGraph()

    def shutdown_tracing():
        """Flush and shut down tracing when the process exits."""
        provider.force_flush()
        provider.shutdown()
        AGENT_GRAPH.close()

    atexit.register(shutdown_tracing)

    return app


if __name__ == "__main__":
    app = create_app()
    app.run(host="0.0.0.0", port=5000, debug=False)
