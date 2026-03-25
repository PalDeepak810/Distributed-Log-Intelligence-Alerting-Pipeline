package monitor

import (
	"strings"
	"sync"
	"time"

	"processor/model"
)

const (
	maxLogs   = 1000
	maxAlerts = 500
)

type AlertEvent struct {
	Type      string `json:"type"`
	Service   string `json:"service"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type Metrics struct {
	TotalLogs      int            `json:"totalLogs"`
	ErrorCount     int            `json:"errorCount"`
	WarnCount      int            `json:"warnCount"`
	LogsPerService map[string]int `json:"logsPerService"`
}

type store struct {
	mu     sync.RWMutex
	logs   []model.Log
	alerts []AlertEvent
}

var global = &store{
	logs:   make([]model.Log, 0, maxLogs),
	alerts: make([]AlertEvent, 0, maxAlerts),
}

func AddLog(log model.Log) {
	global.mu.Lock()
	defer global.mu.Unlock()

	if log.Timestamp == 0 {
		log.Timestamp = time.Now().Unix()
	}

	if len(global.logs) >= maxLogs {
		global.logs = append(global.logs[1:], log)
		return
	}
	global.logs = append(global.logs, log)
}

func AddAlert(event AlertEvent) {
	global.mu.Lock()
	defer global.mu.Unlock()

	if strings.TrimSpace(event.Service) == "" {
		event.Service = "unknown"
	}
	if event.Timestamp == 0 {
		event.Timestamp = time.Now().Unix()
	}

	if len(global.alerts) >= maxAlerts {
		global.alerts = append(global.alerts[1:], event)
		return
	}
	global.alerts = append(global.alerts, event)
}

func GetLogs() []model.Log {
	global.mu.RLock()
	defer global.mu.RUnlock()

	out := make([]model.Log, len(global.logs))
	copy(out, global.logs)
	return out
}

func GetAlerts() []AlertEvent {
	global.mu.RLock()
	defer global.mu.RUnlock()

	out := make([]AlertEvent, len(global.alerts))
	copy(out, global.alerts)
	return out
}

func GetMetrics() Metrics {
	global.mu.RLock()
	defer global.mu.RUnlock()

	metrics := Metrics{
		TotalLogs:      len(global.logs),
		LogsPerService: make(map[string]int),
	}

	for _, log := range global.logs {
		service := strings.TrimSpace(log.Service)
		if service == "" {
			service = "unknown"
		}
		metrics.LogsPerService[service]++

		level := strings.ToUpper(strings.TrimSpace(log.Level))
		if level == "ERROR" {
			metrics.ErrorCount++
		}
		if level == "WARN" {
			metrics.WarnCount++
		}
	}

	return metrics
}
