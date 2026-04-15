/* global document, fetch */

function esc(s) {
  return String(s)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

function checksUrl() {
  const tag = document.getElementById("tag-filter").value.trim();
  return tag ? `/api/v1/checks?tag=${encodeURIComponent(tag)}` : "/api/v1/checks";
}

function exportUrl() {
  const tag = document.getElementById("tag-filter").value.trim();
  return tag
    ? `/api/v1/export/checks.csv?tag=${encodeURIComponent(tag)}`
    : "/api/v1/export/checks.csv";
}

async function parseJSON(r) {
  if (!r.ok) {
    const t = await r.text();
    throw new Error(`HTTP ${r.status}: ${t.slice(0, 200)}`);
  }
  return r.json();
}

async function load() {
  const [sum, checks, drift, fw] = await Promise.all([
    fetch("/api/v1/summary").then(parseJSON),
    fetch(checksUrl()).then(parseJSON),
    fetch("/api/v1/drift").then(parseJSON),
    fetch("/api/v1/frameworks").then(parseJSON),
  ]);

  document.getElementById("summary").innerHTML = `
        <div class="tile"><b>${sum.total}</b> checks</div>
        <div class="tile"><b class="pass">${sum.pass}</b> pass</div>
        <div class="tile"><b class="fail">${sum.fail}</b> fail</div>
        <div class="tile"><b>${sum.drift_count}</b> drift</div>`;

  const tags = fw
    .map((r) => `<span class="tag" title="count">${esc(r.tag)} (${r.count})</span>`)
    .join(" ");
  document.getElementById("frameworks").innerHTML = tags || "<span class='muted'>—</span>";

  const tb = document.getElementById("rows");
  tb.innerHTML = checks
    .map(
      (c) =>
        `<tr><td><code>${esc(c.id)}</code></td><td>${esc(c.title)}</td>` +
        `<td class="${c.expected === "PASS" ? "pass" : "fail"}">${esc(c.expected)}</td>` +
        `<td class="${c.observed === "PASS" ? "pass" : "fail"}">${esc(c.observed)}</td><td>` +
        c.framework_tags.map((t) => `<span class="tag">${esc(t)}</span>`).join("") +
        `</td></tr>`
    )
    .join("");

  document.getElementById("drift").innerHTML = drift.length
    ? drift
        .map(
          (d) =>
            `<tr><td>${esc(d.id)}</td><td>${esc(d.expected)} → ${esc(d.observed)}</td><td>${esc(d.resource)}</td></tr>`
        )
        .join("")
    : `<tr><td colspan="3">No drift.</td></tr>`;
}

function wire() {
  const tag = document.getElementById("tag-filter");
  const exportLink = document.getElementById("export-csv");
  const syncExport = () => exportLink.setAttribute("href", exportUrl());

  tag.addEventListener("change", () => {
    load().catch((e) => console.error(e));
    syncExport();
  });
  tag.addEventListener("input", syncExport);
  document.getElementById("reload").addEventListener("click", () => {
    load().catch((e) => console.error(e));
  });
  syncExport();
}

wire();
load().catch((e) => {
  document.body.innerHTML = "<p>API error: " + esc(e.message) + "</p>";
});
