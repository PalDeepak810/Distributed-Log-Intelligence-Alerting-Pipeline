import { useEffect, useMemo, useRef, useState } from "react";
import "./App.css";

function formatTime(unixSeconds) {
  if (!unixSeconds) return "--";
  return new Date(unixSeconds * 1000).toLocaleTimeString();
}

function formatClock(ms) {
  if (!ms) return "--";
  return new Date(ms).toLocaleTimeString();
}

function App() {
  const [logs, setLogs] = useState([]);
  const [isConnected, setIsConnected] = useState(false);
  const [lastUpdated, setLastUpdated] = useState(0);
  const [levelFilter, setLevelFilter] = useState("ALL");
  const [query, setQuery] = useState("");
  const logsEndRef = useRef(null);

  useEffect(() => {
    const fetchLogs = async () => {
      try {
        const res = await fetch("http://localhost:8090/logs");
        const data = await res.json();
        setLogs(Array.isArray(data) ? data : []);
        setIsConnected(true);
        setLastUpdated(Date.now());
      } catch (err) {
        console.error("Error fetching logs:", err);
        setIsConnected(false);
      }
    };

    fetchLogs();
    const interval = setInterval(fetchLogs, 2000);

    return () => clearInterval(interval);
  }, []);

  const totalLogs = logs.length;
  const errorCount = logs.filter((l) => l.level === "ERROR").length;
  const warnCount = logs.filter((l) => l.level === "WARN").length;

  const serviceCount = {};
  logs.forEach((log) => {
    const service = log.service || "unknown-service";
    serviceCount[service] = (serviceCount[service] || 0) + 1;
  });

  const alerts = logs.filter((l) => l.level === "ERROR");
  const normalizedQuery = query.trim().toLowerCase();

  const filteredLogs = logs.filter((log) => {
    const levelMatch = levelFilter === "ALL" || log.level === levelFilter;
    if (!levelMatch) return false;
    if (!normalizedQuery) return true;

    const haystack = `${log.service || ""} ${log.message || ""}`.toLowerCase();
    return haystack.includes(normalizedQuery);
  });

  const latestLogs = useMemo(() => filteredLogs.slice(-50), [filteredLogs]);

  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [latestLogs.length]);

  const getLevelClass = (level) => {
    if (level === "ERROR") return "level-error";
    if (level === "WARN") return "level-warn";
    return "level-info";
  };

  return (
    <div className="dashboard">
      <header className="hero">
        <div className="hero-top">
          <h1>Log Dashboard</h1>
          <span className={`status-pill ${isConnected ? "up" : "down"}`}>
            {isConnected ? "API Connected" : "API Disconnected"}
          </span>
        </div>
        <p>Real-time visibility into ingestion, alerts, and service behavior.</p>
        <small>Auto refresh: every 2s | Last update: {formatClock(lastUpdated)}</small>
      </header>

      <section className="card">
        <h2>Metrics</h2>
        <div className="metric-grid">
          <div className="metric-box">
            <span>Total Logs</span>
            <strong>{totalLogs}</strong>
          </div>
          <div className="metric-box error">
            <span>Errors</span>
            <strong>{errorCount}</strong>
          </div>
          <div className="metric-box warn">
            <span>Warnings</span>
            <strong>{warnCount}</strong>
          </div>
        </div>

        <div className="service-breakdown">
          <h3>Logs per Service</h3>
          {Object.entries(serviceCount).length === 0 ? (
            <p className="muted">No logs yet.</p>
          ) : (
            Object.entries(serviceCount)
              .sort((a, b) => b[1] - a[1])
              .map(([service, count]) => (
                <div key={service} className="service-row">
                  <span>{service}</span>
                  <strong>{count}</strong>
                </div>
              ))
          )}
        </div>
      </section>

      <section className="card">
        <h2>Alerts</h2>
        {alerts.length === 0 ? (
          <p className="muted">No alerts</p>
        ) : (
          <div className="alert-list">
            {alerts
              .slice()
              .reverse()
              .slice(0, 25)
              .map((log, index) => (
                <div key={`${log.timestamp}-${index}`} className="alert-row">
                  <span className="alert-tag">ERROR</span>
                  <span className="alert-service">{log.service || "unknown-service"}</span>
                  <span className="alert-message">{log.message}</span>
                  <span className="alert-time">{formatTime(log.timestamp)}</span>
                </div>
              ))}
          </div>
        )}
      </section>

      <section className="card">
        <h2>Live Logs</h2>
        <div className="controls">
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search service or message"
            className="search"
          />
          <div className="level-filters">
            {["ALL", "ERROR", "WARN", "INFO"].map((level) => (
              <button
                key={level}
                type="button"
                onClick={() => setLevelFilter(level)}
                className={levelFilter === level ? "active" : ""}
              >
                {level}
              </button>
            ))}
          </div>
        </div>
        <div className="logs-head">
          <span>Service</span>
          <span>Level</span>
          <span>Message</span>
          <span>Time</span>
        </div>
        <div className="logs-body">
          {latestLogs.map((log, index) => (
              <div key={`${log.timestamp}-${index}`} className="log-row">
                <span>{log.service || "unknown-service"}</span>
                <span className={`level-pill ${getLevelClass(log.level)}`}>
                  {log.level || "INFO"}
                </span>
                <span>{log.message || "-"}</span>
                <span>{formatTime(log.timestamp)}</span>
              </div>
            ))}
          {logs.length === 0 ? <p className="muted logs-empty">Waiting for logs...</p> : null}
          {logs.length > 0 && filteredLogs.length === 0 ? (
            <p className="muted logs-empty">No logs match your filter.</p>
          ) : null}
          <div ref={logsEndRef} />
        </div>
      </section>
    </div>
  );
}

export default App;
