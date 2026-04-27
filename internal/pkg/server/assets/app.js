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
            <td><button class="btn btn-primary btn-sm" onclick="resolvePair('${res.RecordA}', '${res.RecordB}', this)">Resolve</button></td>
        `;

    if (res.Resolved) {
        row.style.opacity = '0.5';
        row.querySelector('.status').textContent = 'Resolved';
        row.querySelector('button').disabled = true;
    }
        discoveryBody.prepend(row);
    }

    updateStats();
    setInterval(updateStats, 10000); // Poll every 10 seconds

    window.resolvePair = async (a, b, btn) => {
        try {
            const resp = await fetch('/api/resolve', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({record_a: a, record_b: b})
            });

            if (resp.ok) {
                // Update UI immediately
                const row = btn.closest('tr');
                row.style.transition = 'all 0.5s ease';
                row.style.opacity = '0.3';
                row.style.background = '#f0fff0';
                btn.disabled = true;
                btn.textContent = 'Resolved';
                
                incrementStat(statResolved);
                let potential = parseInt(statDuplicates.textContent);
                statDuplicates.textContent = potential - 1;
            }
        } catch (err) {
            console.error('Resolution failed:', err);
        }
    };
});
