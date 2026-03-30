package pipeline

import (
	"encoding/json"
	"processor/model"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Precompiled regex (performance optimized)
var levelRegex = regexp.MustCompile(`(?i)\b(ERROR|WARN|INFO|DEBUG)\b`)
var cleanRegex = regexp.MustCompile(`(?i)\[(ERROR|WARN|INFO|DEBUG)\]`)

func Normalize(raw map[string]interface{}) model.Log {

	// Structured logs detection
	if raw["service"] != nil || raw["level"] != nil || raw["message"] != nil ||
		raw["status"] != nil || raw["error"] != nil || raw["path"] != nil {
		return normalizeStructured(raw)
	}

	// Raw log handling
	if raw["raw_log"] != nil {
		rawLine := toString(raw["raw_log"])
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

	// Fallback
	return model.Log{
		Service:   "unknown-service",
		Level:     "INFO",
		Message:   "unrecognized format",
		Timestamp: time.Now().Unix(),
	}
}

// STRUCTURED NORMALIZATION

func normalizeStructured(raw map[string]interface{}) model.Log {

	service := getString(raw, "service")
	if service == "" {
		service = getString(raw, "app")
	}
	if service == "" {
		service = getString(raw, "application")
	}
	if service == "" {
		service = "unknown-service"
	}

	level := normalizeLevel(getString(raw, "level"))

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

// RAW LOG NORMALIZATION

func normalizeRaw(log string) model.Log {
	sourceText := extractRawMessage(log)

	level := extractLevel(sourceText)
	cleanMsg := cleanMessage(sourceText)
	if cleanMsg == "" {
		cleanMsg = sourceText
	}

	return model.Log{
		Service:   "unknown-service",
		Level:     level,
		Message:   cleanMsg,
		Timestamp: time.Now().Unix(),
	}
}

func extractRawMessage(log string) string {
	trimmed := strings.TrimSpace(log)
	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(trimmed), &payload); err == nil {
			return firstNonEmpty(
				getString(payload, "raw_log"),
				getString(payload, "message"),
				getString(payload, "error"),
				trimmed,
			)
		}
	}
	return log
}

// HELPERS

func extractLevel(log string) string {
	match := levelRegex.FindString(log)
	if match != "" {
		return strings.ToUpper(match)
	}
	return "INFO"
}

func cleanMessage(log string) string {

	// Remove [ERROR], [WARN], etc.
	clean := cleanRegex.ReplaceAllString(log, "")

	// Remove inline ERROR/WARN/INFO
	clean = levelRegex.ReplaceAllString(clean, "")

	return strings.TrimSpace(clean)
}

func normalizeLevel(level string) string {
	l := strings.ToUpper(strings.TrimSpace(level))

	switch l {
	case "ERROR", "ERR":
		return "ERROR"
	case "WARN", "WARNING":
		return "WARN"
	case "DEBUG":
		return "DEBUG"
	case "INFO":
		return "INFO"
	default:
		return ""
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

		// Try numeric string
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

func toString(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	default:
		data, err := json.Marshal(v)
		if err == nil {
			return string(data)
		}
		return ""
	}
}
