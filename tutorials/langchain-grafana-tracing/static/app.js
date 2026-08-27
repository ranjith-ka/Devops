"use strict";

const state = {
    threadId: localStorage.getItem("tracepilot-thread") || crypto.randomUUID(),
    incidentThreadId: localStorage.getItem("tracepilot-incident-thread") || crypto.randomUUID(),
    busy: false,
};
localStorage.setItem("tracepilot-thread", state.threadId);
localStorage.setItem("tracepilot-incident-thread", state.incidentThreadId);

const $ = (selector, parent = document) => parent.querySelector(selector);
const escapeHtml = (value) => {
    const element = document.createElement("div");
    element.textContent = String(value ?? "");
    return element.innerHTML;
};
const shortId = (value) => value ? `${value.slice(0, 8)}…${value.slice(-6)}` : "—";

async function apiRequest(endpoint, data) {
    const response = await fetch(endpoint, {
        method: data ? "POST" : "GET",
        headers: data ? {"Content-Type": "application/json"} : {},
        body: data ? JSON.stringify(data) : undefined,
    });
    let payload;
    try { payload = await response.json(); }
    catch (_) { payload = {error: `Unexpected response (${response.status})`}; }
    if (!response.ok || payload.success === false) throw new Error(payload.error || "Request failed");
    return payload;
}

function autoSize(textarea) {
    textarea.style.height = "auto";
    textarea.style.height = `${Math.min(textarea.scrollHeight, 180)}px`;
}

function setMode(button, mode) {
    const labels = {question: ["✦", "Ask"], investigate: ["⌁", "Investigate"], compare: ["⇄", "Compare"]};
    button.dataset.mode = mode;
    $(".mode-icon", button).textContent = labels[mode][0];
    $(".mode-label", button).textContent = labels[mode][1];
}

function toast(message) {
    const region = $("#toast-region");
    const item = document.createElement("div");
    item.className = "toast";
    item.textContent = message;
    region.replaceChildren(item);
    window.setTimeout(() => item.remove(), 2400);
}

function showConversation(title = "Trace analysis") {
    $("#welcome").hidden = true;
    $("#conversation").hidden = false;
    $("#conversation-title").textContent = title.length > 54 ? `${title.slice(0, 54)}…` : title;
}

function appendMessage(role, text = "") {
    const article = document.createElement("article");
    article.className = `message ${role}`;
    const avatar = document.createElement("div");
    avatar.className = "avatar";
    avatar.setAttribute("aria-hidden", "true");
    avatar.textContent = role === "user" ? "Y" : "✦";
    const content = document.createElement("div");
    content.className = "message-content";
    const meta = document.createElement("div");
    meta.className = "message-meta";
    meta.textContent = role === "user" ? "You" : "TracePilot";
    const body = document.createElement("div");
    body.className = "message-body";
    body.textContent = text;
    content.append(meta, body);
    article.append(avatar, content);
    $("#timeline").append(article);
    article.scrollIntoView({behavior: "smooth", block: "center"});
    return {article, body, content};
}

function appendLoading(mode) {
    const labels = mode === "investigate"
        ? ["Reading request", "Querying telemetry", "Building evidence"]
        : ["Reading request", "Searching documentation", "Preparing answer"];
    const message = appendMessage("assistant", "");
    message.article.classList.add("loading-message");
    const steps = document.createElement("div");
    steps.className = "loading-steps";
    labels.forEach((label, index) => {
        const span = document.createElement("span");
        span.textContent = label;
        if (index === 0) span.className = "active";
        steps.append(span);
    });
    message.body.replaceChildren(steps);
    let current = 0;
    const timer = window.setInterval(() => {
        steps.children[current].classList.remove("active");
        current = Math.min(current + 1, steps.children.length - 1);
        steps.children[current].classList.add("active");
    }, 1300);
    return {element: message.article, stop: () => window.clearInterval(timer)};
}

function copyButton(value) {
    return `<button class="copy-button" type="button" data-copy="${escapeHtml(value)}" title="Copy trace ID" aria-label="Copy trace ID">Copy</button>`;
}

function rememberTrace(traceId) {
    if (!traceId) return;
    const traces = JSON.parse(localStorage.getItem("tracepilot-recent-traces") || "[]");
    localStorage.setItem("tracepilot-recent-traces", JSON.stringify([traceId, ...traces.filter(id => id !== traceId)].slice(0, 4)));
    renderRecentTraces();
}

function renderRecentTraces() {
    const container = $("#recent-traces");
    const traces = JSON.parse(localStorage.getItem("tracepilot-recent-traces") || "[]");
    if (!traces.length) { container.hidden = true; return; }
    container.hidden = false;
    container.innerHTML = `Recent traces: ${traces.map(id => `<button type="button" data-recent-trace="${escapeHtml(id)}">${escapeHtml(shortId(id))}</button>`).join(" · ")}`;
}

function renderQuestionResult(payload) {
    const message = appendMessage("assistant", payload.answer);
    const meta = document.createElement("div");
    meta.className = "trace-meta";
    meta.innerHTML = `<span class="trace-chip">Trace ${escapeHtml(shortId(payload.trace_id))} ${copyButton(payload.trace_id)}</span>` +
        (payload.sources || []).map(source => `<span class="source-chip">${escapeHtml(source)}</span>`).join("");
    message.content.append(meta);
    rememberTrace(payload.trace_id);
}

function renderInvestigationResult(payload) {
    const message = appendMessage("assistant", payload.answer);
    const meta = document.createElement("div");
    meta.className = "trace-meta";
    meta.innerHTML = `<span class="trace-chip">Investigation ${escapeHtml(shortId(payload.trace_id))} ${copyButton(payload.trace_id)}</span>` +
        (payload.selected_tools || []).map(tool => `<span class="tool-chip">${escapeHtml(tool)}</span>`).join("");
    message.content.append(meta);
    if (payload.evidence && payload.evidence.length) {
        const grid = document.createElement("div");
        grid.className = "evidence-grid";
        payload.evidence.forEach(item => {
            const card = document.createElement("section");
            card.className = "evidence-card";
            const heading = document.createElement("h4");
            heading.textContent = item.tool;
            const pre = document.createElement("pre");
            pre.textContent = JSON.stringify(item.data, null, 2);
            card.append(heading, pre);
            grid.append(card);
        });
        message.content.append(grid);
    }
    rememberTrace(payload.trace_id);
}

function renderWaterfall(trace, label, sharedDuration) {
    const rows = (trace.spans || []).map(span => {
        const left = Math.min(100, span.start_offset_ms / sharedDuration * 100);
        const width = Math.max(.7, Math.min(100 - left, span.duration_ms / sharedDuration * 100));
        const classes = ["waterfall-bar", span.is_root ? "root" : "", span.status === "error" ? "error" : ""].join(" ");
        const tip = `${span.name}: ${Number(span.duration_ms).toFixed(2)} ms · ${span.status}`;
        return `<div class="waterfall-row"><div class="waterfall-label" title="${escapeHtml(span.name)}">${escapeHtml(span.name)}</div><div class="waterfall-track"><div class="${classes}" style="left:${left}%;width:${width}%" title="${escapeHtml(tip)}"></div></div></div>`;
    }).join("");
    return `<section class="waterfall-panel"><h4 class="waterfall-title">${label} · ${Number(trace.duration_ms).toFixed(2)} ms</h4><div class="waterfall-id">${escapeHtml(trace.trace_id)}</div>${rows}</section>`;
}

function renderLogs(logs, label) {
    const rows = !logs.length ? `<p class="empty-state">No correlated logs found</p>` : logs.slice(0, 50).map(log => {
        const time = new Date(Number(log.timestamp) / 1e6).toLocaleTimeString();
        const detail = [log.node, log.error_type, log.tool].filter(Boolean).join(" · ");
        return `<div class="log-entry"><span>${escapeHtml(time)}</span><span><span class="log-event">${escapeHtml(log.event || "log")}</span>${detail ? ` · ${escapeHtml(detail)}` : ""}</span></div>`;
    }).join("");
    return `<section class="log-panel"><h4>${label}</h4>${rows}</section>`;
}

function renderComparison(payload) {
    const comparison = payload.comparison;
    const diff = Number(comparison.root_diff_ms);
    const slower = diff >= 0;
    const sharedDuration = Math.max(comparison.trace_a_waterfall.duration_ms, comparison.trace_b_waterfall.duration_ms, 1);
    const spanRows = Object.entries(comparison.span_details || {}).sort((a, b) => Math.abs(b[1].delta_ms) - Math.abs(a[1].delta_ms)).map(([name, span]) => {
        const delta = Number(span.delta_ms);
        return `<tr><td>${escapeHtml(name)}</td><td>${Number(span.trace_a_ms).toFixed(2)} ms</td><td>${Number(span.trace_b_ms).toFixed(2)} ms</td><td class="${delta > 0 ? "delta-bad" : "delta-good"}">${delta > 0 ? "+" : ""}${delta.toFixed(2)} ms</td><td>${escapeHtml(span.status_a)} / ${escapeHtml(span.status_b)}</td></tr>`;
    }).join("");
    const message = appendMessage("assistant", "Comparison complete. Here is the strongest signal from the two traces.");
    const card = document.createElement("section");
    card.className = "comparison-card";
    card.innerHTML = `<div class="summary-grid"><div class="metric"><small>Baseline</small><strong>${Number(comparison.trace_a_waterfall.duration_ms).toFixed(2)} ms</strong></div><div class="metric"><small>Candidate</small><strong>${Number(comparison.trace_b_waterfall.duration_ms).toFixed(2)} ms</strong></div><div class="metric"><small>Change</small><strong class="${slower ? "bad" : "good"}">${slower ? "+" : ""}${diff.toFixed(2)} ms</strong></div></div><div class="finding"><strong>Root-cause hypothesis</strong><p>${escapeHtml(comparison.root_cause)}</p></div><div class="finding"><strong>Recommended next action</strong><p>${escapeHtml(comparison.recommendation)}</p></div><div class="waterfall-grid">${renderWaterfall(comparison.trace_a_waterfall, "Baseline A", sharedDuration)}${renderWaterfall(comparison.trace_b_waterfall, "Candidate B", sharedDuration)}</div><details class="details-toggle"><summary>View correlated logs</summary><div class="log-grid">${renderLogs(comparison.trace_a_logs || [], "Baseline logs")}${renderLogs(comparison.trace_b_logs || [], "Candidate logs")}</div></details>${spanRows ? `<details class="details-toggle"><summary>View span deltas</summary><div class="table-wrap"><table><thead><tr><th>Span</th><th>Baseline</th><th>Candidate</th><th>Delta</th><th>Status A/B</th></tr></thead><tbody>${spanRows}</tbody></table></div></details>` : ""}`;
    message.content.append(card);
    rememberTrace(comparison.trace_a_waterfall.trace_id);
    rememberTrace(comparison.trace_b_waterfall.trace_id);
}

function showError(error) {
    const message = appendMessage("assistant", "");
    const card = document.createElement("div");
    card.className = "error-card";
    card.textContent = error.message || "Something went wrong. Check the local services and try again.";
    message.body.replaceChildren(card);
}

async function submitPrompt(input, modeButton) {
    if (state.busy) return;
    const question = input.value.trim();
    const mode = modeButton.dataset.mode;
    if (mode === "compare") { openCompareForm(); return; }
    if (!question) { input.focus(); toast("Enter a question first"); return; }
    state.busy = true;
    showConversation(question);
    appendMessage("user", question);
    input.value = "";
    autoSize(input);
    const loading = appendLoading(mode);
    document.querySelectorAll(".send-button").forEach(button => button.disabled = true);
    try {
        if (mode === "investigate") {
            const payload = await apiRequest("/api/investigate", {question, thread_id: state.incidentThreadId, max_model_calls: 6});
            loading.stop(); loading.element.remove(); renderInvestigationResult(payload);
        } else {
            const payload = await apiRequest("/api/question", {question, thread_id: state.threadId});
            loading.stop(); loading.element.remove(); renderQuestionResult(payload);
        }
    } catch (error) { loading.stop(); loading.element.remove(); showError(error); }
    finally { state.busy = false; document.querySelectorAll(".send-button").forEach(button => button.disabled = false); $("#dock-input").focus(); }
}

function openCompareForm() {
    showConversation("Compare two traces");
    const message = appendMessage("assistant", "Paste a baseline and candidate trace ID. I’ll compare timing, structure, status, and correlated logs.");
    const form = $("#compare-form-template").content.firstElementChild.cloneNode(true);
    message.content.append(form);
    const recent = JSON.parse(localStorage.getItem("tracepilot-recent-traces") || "[]");
    if (recent[0]) form.elements.trace_a.value = recent[0];
    if (recent[1]) form.elements.trace_b.value = recent[1];
    form.addEventListener("submit", async event => {
        event.preventDefault();
        const traceA = form.elements.trace_a.value.trim();
        const traceB = form.elements.trace_b.value.trim();
        if (!/^[0-9a-f]{32}$/i.test(traceA) || !/^[0-9a-f]{32}$/i.test(traceB)) { toast("Each trace ID must be 32 hexadecimal characters"); return; }
        form.querySelector("button").disabled = true;
        const loading = appendLoading("investigate");
        try { const payload = await apiRequest("/api/compare", {trace_a: traceA, trace_b: traceB}); loading.stop(); loading.element.remove(); renderComparison(payload); }
        catch (error) { loading.stop(); loading.element.remove(); showError(error); }
        finally { form.querySelector("button").disabled = false; }
    });
    form.elements.trace_a.focus();
}

async function loadStatus() {
    try {
        const payload = await apiRequest("/api/status");
        const available = payload.services.filter(item => item.available).length;
        const dot = $(".status-dot", $("#service-status"));
        dot.className = `status-dot ${available === payload.services.length ? "" : available ? "partial" : "offline"}`;
        $(".status-summary span:last-child").textContent = available === payload.services.length ? "All services connected" : `${available}/${payload.services.length} services connected`;
        $("#status-popover").innerHTML = payload.services.map(item => `<div class="status-row"><span>${escapeHtml(item.name)}</span><b class="${item.available ? "" : "off"}">${item.available ? "Connected" : "Offline"}</b></div>`).join("");
    } catch (_) {
        $(".status-dot", $("#service-status")).className = "status-dot offline";
        $(".status-summary span:last-child").textContent = "Status unavailable";
    }
}

function resetSession() {
    state.threadId = crypto.randomUUID(); state.incidentThreadId = crypto.randomUUID();
    localStorage.setItem("tracepilot-thread", state.threadId); localStorage.setItem("tracepilot-incident-thread", state.incidentThreadId);
    $("#timeline").replaceChildren(); $("#conversation").hidden = true; $("#welcome").hidden = false; $("#hero-input").focus();
}

document.addEventListener("DOMContentLoaded", () => {
    const composers = [["#hero-input", "#hero-send", "#hero-mode"], ["#dock-input", "#dock-send", "#dock-mode"]];
    composers.forEach(([inputSelector, sendSelector, modeSelector]) => {
        const input = $(inputSelector), send = $(sendSelector), mode = $(modeSelector), menu = $(`${modeSelector}-menu`);
        input.addEventListener("input", () => autoSize(input));
        input.addEventListener("keydown", event => { if (event.key === "Enter" && !event.shiftKey) { event.preventDefault(); submitPrompt(input, mode); } });
        send.addEventListener("click", () => submitPrompt(input, mode));
        mode.addEventListener("click", event => { event.stopPropagation(); document.querySelectorAll(".mode-menu").forEach(item => { if (item !== menu) item.classList.remove("open"); }); menu.classList.toggle("open"); });
        menu.addEventListener("click", event => { const choice = event.target.closest("[data-select-mode]"); if (!choice) return; setMode(mode, choice.dataset.selectMode); menu.classList.remove("open"); if (choice.dataset.selectMode === "compare") openCompareForm(); else input.focus(); });
    });
    document.addEventListener("click", event => {
        if (!event.target.closest(".mode-button, .mode-menu")) document.querySelectorAll(".mode-menu").forEach(menu => menu.classList.remove("open"));
        const copy = event.target.closest("[data-copy]"); if (copy) navigator.clipboard.writeText(copy.dataset.copy).then(() => toast("Trace ID copied"));
        const recent = event.target.closest("[data-recent-trace]"); if (recent) { setMode($("#hero-mode"), "investigate"); $("#hero-input").value = `Investigate trace ${recent.dataset.recentTrace}`; $("#hero-input").focus(); }
    });
    document.querySelectorAll(".suggestion-card").forEach(card => card.addEventListener("click", () => { const mode = card.dataset.mode; setMode($("#hero-mode"), mode); if (mode === "compare") openCompareForm(); else { $("#hero-input").value = card.dataset.prompt || ""; autoSize($("#hero-input")); $("#hero-input").focus(); } }));
    $("#new-session").addEventListener("click", resetSession); $("#clear-conversation").addEventListener("click", resetSession);
    renderRecentTraces(); loadStatus(); window.setInterval(loadStatus, 30000);
});
