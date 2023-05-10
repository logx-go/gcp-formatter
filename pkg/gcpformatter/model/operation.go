package model

// Operation according to https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntryOperation
type Operation struct {
	Id       string `json:"id,omitempty"`       // Optional. An arbitrary operation identifier. Log entries with the same identifier are assumed to be part of the same operation.
	Producer string `json:"producer,omitempty"` // Optional. An arbitrary producer identifier. The combination of id and producer must be globally unique. Examples for producer: "MyDivision.MyBigCompany.com", "github.com/MyProject/MyApplication".
	First    bool   `json:"first,omitempty"`    // Optional. Set this to True if this is the first log entry in the operation.
	Last     bool   `json:"last,omitempty"`     // Optional. Set this to True if this is the last log entry in the operation.
}
