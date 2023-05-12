package model

import (
	"encoding/json"
	"time"
)

// LogEntry according to https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
type LogEntry struct {
	Severity       string                     `json:"severity,omitempty"`
	InsertId       string                     `json:"insertId,omitempty"`
	TraceSampled   bool                       `json:"traceSampled,omitempty"`
	Trace          string                     `json:"trace,omitempty"`
	SpanId         string                     `json:"spanId,omitempty"`
	JsonPayload    map[string]json.RawMessage `json:"jsonPayload,omitempty"`
	Labels         map[string]string          `json:"labels,omitempty"`
	HttpRequest    *HttpRequest               `json:"httpRequest,omitempty"`
	Operation      *Operation                 `json:"operation,omitempty"`
	SourceLocation *SourceLocation            `json:"sourceLocation,omitempty"`
	Timestamp      time.Time                  `json:"timestamp,omitempty"`
}
