function app() {
  return {
    // Connection
    source: 'kafka',
    destType: 'topic',
    name: '',
    isConnected: false,
    eventSource: null,
    errorMsg: '',

    // View
    viewMode: 'text',
    filter: '',
    fromDate: '',
    columnFilters: {},

    // Events
    events: [],
    autoScroll: true,

    init() {
      // Default "from" to current date/time
      this.fromDate = toDatetimeLocal(new Date());

      // Restore last used settings from localStorage
      const saved = localStorage.getItem('consumer-settings');
      if (saved) {
        const s = JSON.parse(saved);
        this.source = s.source || 'kafka';
        this.destType = s.destType || 'topic';
        this.name = s.name || '';
        this.viewMode = s.viewMode || 'text';
      }
    },

    // --- Computed ---

    get tableColumns() {
      const fixed = ['time', 'event_source', 'topic/queue'];
      const fixedSet = new Set(fixed);
      const dynamic = new Set();
      for (const event of this.events) {
        try {
          const obj = JSON.parse(event.Payload);
          if (obj && typeof obj === 'object' && !Array.isArray(obj)) {
            Object.keys(obj).forEach(k => { if (!fixedSet.has(k)) dynamic.add(k); });
          }
        } catch { /* not JSON */ }
      }
      return [...fixed, ...dynamic];
    },

    get filteredEvents() {
      const fromMs = this.fromDate ? new Date(this.fromDate).getTime() : 0;

      if (this.viewMode === 'table') {
        return this.events.filter(event => {
          if (fromMs && this.eventTimeMs(event) < fromMs) return false;
          for (const [col, val] of Object.entries(this.columnFilters)) {
            if (!val) continue;
            const needle = val.toLowerCase();
            const cell = String(this.getCell(event, col)).toLowerCase();
            if (!cell.includes(needle)) return false;
          }
          return true;
        });
      }

      const needle = this.filter ? this.filter.toLowerCase() : null;
      return this.events.filter(e => {
        if (fromMs && this.eventTimeMs(e) < fromMs) return false;
        if (!needle) return true;
        return (
          e.Payload.toLowerCase().includes(needle) ||
          e.Source.toLowerCase().includes(needle) ||
          e.Name.toLowerCase().includes(needle)
        );
      });
    },

    // --- Actions ---

    start() {
      this.errorMsg = '';
      this.saveSettings();
      this.events = [];

      const params = new URLSearchParams({
        source: this.source,
        name:   this.name.trim(),
        type:   this.source === 'kafka' ? 'topic' : this.destType,
      });

      if (this.eventSource) this.eventSource.close();
      this.eventSource = new EventSource('/consumer/events?' + params.toString());

      this.eventSource.onmessage = (e) => {
        const event = JSON.parse(e.data);
        this.events.push(event);
        if (this.autoScroll) {
          this.$nextTick(() => {
            const list = document.getElementById('events-list');
            if (list) list.scrollTop = list.scrollHeight;
          });
        }
      };

      this.eventSource.addEventListener('error', (e) => {
        // Only treat it as a fatal error if the connection was never established
        // or the server sent an explicit error event.
        if (e.type === 'error' && this.eventSource.readyState === EventSource.CLOSED) {
          this.isConnected = false;
          this.eventSource = null;
          this.errorMsg = 'Connection lost.';
        }
      });

      this.isConnected = true;
    },

    stop() {
      if (this.eventSource) {
        this.eventSource.close();
        this.eventSource = null;
      }
      this.isConnected = false;
    },

    clear() {
      this.events = [];
    },

    // --- Helpers ---

    formatTime(ts) {
      if (!ts) return '';
      const d = new Date(ts);
      return d.toLocaleTimeString('en-US', {
        hour12: false,
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
      });
    },

    highlight(payload) {
      try {
        const obj = JSON.parse(payload);
        return syntaxHighlight(JSON.stringify(obj, null, 2));
      } catch {
        return escapeHtml(payload);
      }
    },

    // Returns the event's effective timestamp in ms:
    // uses payload.created_at (epoch seconds) if present, else ReceivedAt.
    eventTimeMs(event) {
      try {
        const obj = JSON.parse(event.Payload);
        if (obj && typeof obj.created_at === 'number') return obj.created_at * 1000;
      } catch {}
      return new Date(event.ReceivedAt).getTime();
    },

    getCell(event, col) {
      if (col === 'time')         return this.formatTime(event.ReceivedAt);
      if (col === 'event_source') return event.Source;
      if (col === 'topic/queue')  return event.Name;
      try {
        const obj = JSON.parse(event.Payload);
        const val = obj[col];
        if (val === undefined) return '';
        // Format epoch-seconds timestamps as readable dates
        if (col === 'created_at' && typeof val === 'number') {
          return formatEpochSeconds(val);
        }
        return typeof val === 'object' ? JSON.stringify(val) : String(val);
      } catch {
        return '';
      }
    },

    tableHeader() {
      return this.tableColumns.map(col =>
        `<th><div class="col-label">${escapeHtml(col)}</div><input type="text" class="col-filter" data-col="${escapeHtml(col)}" placeholder="filter…"></th>`
      ).join('');
    },

    tableBody() {
      return this.filteredEvents.map((event, i) => {
        const cells = this.tableColumns.map(col =>
          `<td>${escapeHtml(String(this.getCell(event, col)))}</td>`
        ).join('');
        return `<tr class="${i % 2 === 0 ? 'row-even' : 'row-odd'}">${cells}</tr>`;
      }).join('');
    },

    saveSettings() {
      localStorage.setItem('consumer-settings', JSON.stringify({
        source: this.source,
        destType: this.destType,
        name: this.name,
        viewMode: this.viewMode,
      }));
    },
  };
}

function syntaxHighlight(json) {
  const escaped = escapeHtml(json);
  return escaped.replace(
    /("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g,
    match => {
      let cls = 'jn'; // number
      if (/^"/.test(match)) {
        cls = /:$/.test(match) ? 'jk' : 'js'; // key or string
      } else if (/true|false/.test(match)) {
        cls = 'jb'; // bool
      } else if (/null/.test(match)) {
        cls = 'jnull';
      }
      return `<span class="${cls}">${match}</span>`;
    }
  );
}

// Returns a string like "2024-03-10 19:06:10" from epoch seconds
function formatEpochSeconds(epochSeconds) {
  const d = new Date(epochSeconds * 1000);
  const pad = n => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ` +
         `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}

// Formats a Date to "YYYY-MM-DDTHH:MM" for datetime-local inputs
function toDatetimeLocal(date) {
  const pad = n => String(n).padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth()+1)}-${pad(date.getDate())}` +
         `T${pad(date.getHours())}:${pad(date.getMinutes())}`;
}

function escapeHtml(str) {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}
