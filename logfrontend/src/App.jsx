import { useEffect, useMemo, useState } from "react";

const API_BASE = "http://localhost:8090";

function formatTime(unixSeconds) {
  if (!unixSeconds) return "-";
  return new Date(unixSeconds * 1000).toLocaleTimeString();
}

function App() {
  const [logs, setLogs] = useState([]);
  const [alerts, setAlerts] = useState([]);
  const [metrics, setMetrics] = useState({
    totalLogs: 0,
    errorCount: 0,
    warnCount: 0,
    logsPerService: {}
  });
  const [error, setError] = useState("");

  useEffect(() => {
    let active = true;

    const fetchData = async () => {
      try {
        const [logsRes, alertsRes, metricsRes] = await Promise.all([
          fetch(`${API_BASE}/logs`),
          fetch(`${API_BASE}/alerts`),
          fetch(`${API_BASE}/metrics`)
        ]);

        if (!logsRes.ok || !alertsRes.ok || !metricsRes.ok) {
          throw new Error("Failed to fetch monitoring data");
        }

        const [logsData, alertsData, metricsData] = await Promise.all([
          logsRes.json(),
          alertsRes.json(),
          metricsRes.json()
        ]);

        if (!active) return;

        setLogs(Array.isArray(logsData) ? logsData : []);
        setAlerts(Array.isArray(alertsData) ? alertsData : []);
        setMetrics(metricsData || {});
        setError("");
      } catch (e) {
        if (!active) return;
        setError("Could not connect to processor API at http://localhost:8090");
      }
    };

    fetchData();
    const intervalId = setInterval(fetchData, 2000);

    return () => {
      active = false;
      clearInterval(intervalId);
    };
  }, []);

  const importantAlerts = useMemo(() => {
    return alerts.filter(
      (alert) =>
        alert.type === "ERROR_LOG" || alert.type === "SPIKE_DETECTED"
    );
  }, [alerts]);

  const serviceEntries = useMemo(() => {
    return Object.entries(metrics.logsPerService || {}).sort((a, b) => b[1] - a[1]);
  }, [metrics.logsPerService]);

  return (
    <div className="page">
      <header className="header">
        <h1>Log Frontend</h1>
        <p>Real-time pipeline visibility for logs, alerts, and metrics.</p>
      </header>

      {error ? <div className="errorBanner">{error}</div> : null}

      <section className="panel">
        <h2>Metrics Summary</h2>
        <div className="summaryGrid">
          <div className="card">
            <span className="label">Total Logs</span>
            <span className="value">{metrics.totalLogs ?? 0}</span>
          </div>
          <div className="card">
            <span className="label">ERROR Count</span>
            <span className="value error">{metrics.errorCount ?? 0}</span>
          </div>
          <div className="card">
            <span className="label">WARN Count</span>
            <span className="value warn">{metrics.warnCount ?? 0}</span>
          </div>
        </div>

        <div className="serviceList">
          <h3>Logs Per Service</h3>
          {serviceEntries.length === 0 ? (
            <p>No logs received yet.</p>
          ) : (
            serviceEntries.map(([service, count]) => (
              <div className="serviceRow" key={service}>
                <span>{service}</span>
                <strong>{count}</strong>
              </div>
            ))
          )}
        </div>
      </section>

      <section className="panel">
        <h2>Alerts Panel</h2>
        {importantAlerts.length === 0 ? (
          <p>No active alerts yet.</p>
        ) : (
          <div className="alerts">
            {importantAlerts
              .slice()
              .reverse()
              .slice(0, 30)
              .map((alert, idx) => (
                <div className="alertRow" key={`${alert.timestamp}-${idx}`}>
                  <span className="alertType">{alert.type}</span>
                  <span>{alert.service || "unknown"}</span>
                  <span>{alert.message}</span>
                  <span>{formatTime(alert.timestamp)}</span>
                </div>
              ))}
          </div>
        )}
      </section>

      <section className="panel">
        <h2>Live Logs Stream</h2>
        <div className="tableHead">
          <span>Service</span>
          <span>Level</span>
          <span>Message</span>
          <span>Time</span>
        </div>
        <div className="tableBody">
          {logs
            .slice()
            .reverse()
            .slice(0, 100)
            .map((log, index) => (
              <div className="tableRow" key={`${log.timestamp}-${index}`}>
                <span>{log.service || "unknown"}</span>
                <span className={`pill ${String(log.level || "").toLowerCase()}`}>
                  {log.level || "INFO"}
                </span>
                <span>{log.message || "-"}</span>
                <span>{formatTime(log.timestamp)}</span>
              </div>
            ))}
        </div>
      </section>
    </div>
  );
}

export default App;
