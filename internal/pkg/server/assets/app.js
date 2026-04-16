document.addEventListener('DOMContentLoaded', () => {
    const discoveryBody = document.getElementById('discovery-body');
    const statScanned = document.getElementById('stat-scanned');
    const statDuplicates = document.getElementById('stat-duplicates');
    const statResolved = document.getElementById('stat-resolved');

    // 1. WebSocket Setup
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${protocol}//${window.location.host}/ws`);

    ws.onmessage = (event) => {
        const result = JSON.parse(event.data);
        addDiscoveryRow(result);
        incrementStat(statDuplicates);
    };

    ws.onclose = () => {
        console.log('WebSocket connection closed. Attempting reconnect...');
        setTimeout(() => window.location.reload(), 5000);
    };

    // 2. Statistics Polling
    async function updateStats() {
        try {
            const resp = await fetch('/api/stats');
            const data = await resp.json();
            statScanned.textContent = data.total_records;
            statDuplicates.textContent = data.duplicates;
            statResolved.textContent = data.resolved;
        } catch (err) {
            console.error('Failed to fetch stats:', err);
        }
    }

    function incrementStat(element) {
        let val = parseInt(element.textContent);
        element.textContent = val + 1;
    }

    function addDiscoveryRow(res) {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${res.RecordA} <span class="text-muted">↔</span> ${res.RecordB}</td>
            <td><span class="similarity-badge">${(res.Score * 100).toFixed(2)}%</span></td>
            <td>${res.Algorithm}</td>
            <td><span class="status">Potential</span></td>
            <td><button class="btn btn-primary btn-sm" onclick="alert('Resolution not implemented in prototype')">Resolve</button></td>
        `;
        discoveryBody.prepend(row);
    }

    // Initial load
    updateStats();
    setInterval(updateStats, 10000); // Poll every 10 seconds
});
