const API_BASE = '/api/v1';

async function api(path, opts = {}) {
  const res = await fetch(`${API_BASE}${path}`, opts);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

function escape(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

// Tab switching
document.querySelectorAll('.tab').forEach(tab => {
  tab.addEventListener('click', () => {
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    document.querySelectorAll('.tab-panel').forEach(p => p.classList.remove('active'));
    tab.classList.add('active');
    document.getElementById(tab.dataset.tab).classList.add('active');
  });
});

// Upload document
document.getElementById('uploadDoc')?.addEventListener('click', async () => {
  const file = document.getElementById('docFile').files[0];
  if (!file) return alert('Select a file');
  const fd = new FormData();
  fd.append('file', file);
  const data = await api('/documents', { method: 'POST', body: fd });
  document.getElementById('docResult').innerHTML =
    `<div class="badge green">${data.claims_found} claims found</div>`;
  renderClaimList(data.claims);
});

// Upload evidence
document.getElementById('uploadEv')?.addEventListener('click', async () => {
  const file = document.getElementById('evFile').files[0];
  if (!file) return alert('Select a file');
  const fd = new FormData();
  fd.append('file', file);
  const data = await api('/evidence', { method: 'POST', body: fd });
  document.getElementById('evResult').innerHTML =
    `<div class="badge green">${data.evidence.source_type} loaded</div>`;
});

// Run analysis
document.getElementById('runAnalysis')?.addEventListener('click', async () => {
  const docs = document.getElementById('analysisDocs').value.split('\n').filter(Boolean);
  const ev = document.getElementById('analysisEv').value.split('\n').filter(Boolean);
  const fd = new FormData();
  docs.forEach(d => fd.append('documents', d));
  ev.forEach(e => fd.append('evidence', e));
  const data = await api('/analyze', { method: 'POST', body: fd });
  renderSummary(data.summary);
  renderFindings(data.assumptions, data.verifications);
  renderGaps(data.gaps);
  renderGraph();
});

// Load results
async function loadResults() {
  const data = await api('/summary');
  renderSummary(data.summary);
  await renderGraph();
}

function renderSummary(s) {
  const el = document.getElementById('summaryContent');
  if (!el) return;
  el.innerHTML = `
    <div class="stat"><span class="stat-label">Claims</span><span class="stat-value">${s.claims_found}</span></div>
    <div class="stat"><span class="stat-label">Assumptions</span><span class="stat-value">${s.assumptions}</span></div>
    <div class="stat stat-green"><span class="stat-label">Verified</span><span class="stat-value">${s.verified}</span></div>
    <div class="stat stat-red"><span class="stat-label">Contradicted</span><span class="stat-value">${s.contradicted}</span></div>
    <div class="stat stat-yellow"><span class="stat-label">Unknown</span><span class="stat-value">${s.unknown}</span></div>
    <div class="stat stat-critical"><span class="stat-label">Critical Gaps</span><span class="stat-value">${s.critical_gaps}</span></div>
  `;
  document.getElementById('tab-summary').click();
}

function renderClaimList(claims) {
  const el = document.getElementById('claimList');
  if (!el) return;
  el.innerHTML = claims.map(c => `<div class="finding-item"><strong>${escape(c.text)}</strong> <span class="badge">${Math.round(c.extraction_confidence * 100)}%</span></div>`).join('');
}

function renderFindings(assumptions, verifications) {
  const el = document.getElementById('findingsContent');
  if (!el) return;
  const verMap = Object.fromEntries(verifications.map(v => [v.assumption_id, v]));

  el.innerHTML = assumptions.map(a => {
    const v = verMap[a.id];
    const status = v ? v.result : 'PENDING';
    const conf = v ? Math.round(v.confidence * 100) : 0;
    const cls = status === 'VERIFIED' ? 'green' : status === 'CONTRADICTED' ? 'red' : 'yellow';
    return `<div class="finding-item">
      <div class="finding-status ${cls}">${status}</div>
      <div><strong>${escape(a.text)}</strong></div>
      <div class="finding-meta">Confidence: ${conf}% | ${a.assumption_type} | Evidence: ${v ? v.evidence_used.length : 0}</div>
      ${v && v.reasoning ? `<div class="finding-meta">${escape(v.reasoning)}</div>` : ''}
    </div>`;
  }).join('');
  document.getElementById('tab-findings').click();
}

function renderGaps(gaps) {
  const el = document.getElementById('gapsContent');
  if (!el) return;
  el.innerHTML = gaps.map(g => {
    const sev = g.severity.toLowerCase();
    return `<div class="finding-item gap-${sev}">
      <div class="finding-status ${sev === 'critical' ? 'red' : sev === 'high' ? 'orange1' : 'yellow'}">${g.severity}</div>
      <div><strong>${g.type}</strong></div>
      <div>${escape(g.description)}</div>
    </div>`;
  }).join('');
  document.getElementById('tab-gaps').click();
}

async function renderGraph() {
  const el = document.getElementById('graphContent');
  if (!el) return;
  try {
    const data = await api('/graph');
    const nodeTypes = {};
    data.nodes.forEach(n => {
      nodeTypes[n.node_type] = (nodeTypes[n.node_type] || 0) + 1;
    });
    const relTypes = {};
    data.edges.forEach(e => {
      relTypes[e.relationship] = (relTypes[e.relationship] || 0) + 1;
    });
    el.innerHTML = `
      <div class="graph-info">
        <div class="stat"><span class="stat-label">Total Nodes</span><span class="stat-value">${data.node_count}</span></div>
        <div class="stat"><span class="stat-label">Total Edges</span><span class="stat-value">${data.edge_count}</span></div>
      </div>
      <h3>Node Types</h3>
      <div class="graph-info">${Object.entries(nodeTypes).map(([k, v]) => `<div class="stat"><span class="stat-label">${k}</span><span class="stat-value">${v}</span></div>`).join('')}</div>
      <h3>Relationships</h3>
      <div class="graph-info">${Object.entries(relTypes).map(([k, v]) => `<div class="stat"><span class="stat-label">${k}</span><span class="stat-value">${v}</span></div>`).join('')}</div>
      <pre class="graph-json">${escape(JSON.stringify(data, null, 2).slice(0, 2000))}...</pre>
    `;
  } catch (e) {
    el.innerHTML = `<div class="finding-item">Run an analysis first to see the graph.</div>`;
  }
}

// Initial load
document.addEventListener('DOMContentLoaded', () => {
  renderGraph();
});
