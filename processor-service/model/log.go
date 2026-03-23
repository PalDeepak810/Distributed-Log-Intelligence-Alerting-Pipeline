package model

type Log struct {
	Service   string `json:"service"`
	Instance  string `json:"instance"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	TraceID   string `json:"traceId"`
}