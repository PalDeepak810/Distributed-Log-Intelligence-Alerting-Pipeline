package pipeline

import (
	"fmt"
	"processor/model"
	"strconv"
	"strings"
	"time"
)

func normalizeStructured(raw map[string]interface{}) model.Log {
	service := getString(raw, "service")
	if service == "" {
		service = getString(raw, "app")
	}
	if service == "" {
		service = getString(raw, "application")
	}
	if service == "" {
		service = "unknown"
	}

	level := strings.ToUpper(strings.TrimSpace(getString(raw, "level")))
	status := getInt(raw, "status")
	errField := getString(raw, "error")
	if level == "" {
		switch {
		case errField != "":
			level = "ERROR"
		case status >= 500:
			level = "ERROR"
		case status >= 400:
			level = "WARN"
		default:
			level = "INFO"
		}
	}

	message := getString(raw, "message")
	if message == "" {
		message = errField
	}
	if message == "" {
		message = getString(raw, "path")
	}

	traceID := getString(raw, "traceId")
	if traceID == "" {
		traceID = getString(raw, "trace_id")
	}

	return model.Log{
		Service:   service,
		Instance:  firstNonEmpty(getString(raw, "instance"), getString(raw, "host")),
		Level:     level,
		Message:   message,
		Timestamp: getTimestamp(raw),
		TraceID:   traceID,
	}
}

func normalizeRaw(log string) model.Log {

	lower := strings.ToLower(log)

	level := "INFO"

	if strings.Contains(lower, "error") {
		level = "ERROR"
	} else if strings.Contains(lower, "warn") {
		level = "WARN"
	}

	return model.Log{
		Service:   "unknown",
		Level:     level,
		Message:   log,
		Timestamp: time.Now().Unix(),
	}
}

func getString(raw map[string]interface{}, key string) string {

	if val, ok := raw[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getInt(raw map[string]interface{}, key string) int {
	val, ok := raw[key]
	if !ok || val == nil {
		return 0
	}

	switch v := val.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err == nil {
			return n
		}
	}

	return 0
}

func getTimestamp(raw map[string]interface{}) int64 {
	val, ok := raw["timestamp"]
	if !ok || val == nil {
		return time.Now().Unix()
	}

	switch v := val.(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return time.Now().Unix()
		}
		if parsed, err := time.Parse(time.RFC3339Nano, s); err == nil {
			return parsed.Unix()
		}
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			return n
		}
	}

	return time.Now().Unix()
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func Normalize(raw map[string]interface{}) model.Log {
	// structured logs
	if raw["service"] != nil || raw["level"] != nil || raw["message"] != nil ||
		raw["status"] != nil || raw["error"] != nil || raw["path"] != nil {
		return normalizeStructured(raw)
	}

	// unstructured raw logs
	if raw["raw_log"] != nil {
		rawLine := getString(raw, "raw_log")
		if rawLine == "" {
			rawLine = fmt.Sprintf("%v", raw["raw_log"])
		}

		normalized := normalizeRaw(rawLine)
		normalized.Service = firstNonEmpty(
			getString(raw, "service"),
			getString(raw, "source"),
			getString(raw, "app"),
			getString(raw, "application"),
			normalized.Service,
		)
		normalized.Timestamp = getTimestamp(raw)
		return normalized
	}

	// fallback
	return model.Log{
		Service:   "unknown",
		Level:     "INFO",
		Message:   "unrecognized format",
		Timestamp: time.Now().Unix(),
	}
}
