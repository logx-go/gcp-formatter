package gcp_formatter

// special field names
const (
	FieldNameInsertId                       string = "gcp:insert_id"
	FieldNameOperationId                    string = "gcp:operation_id"
	FieldNameOperationProducer              string = "gcp:operation_producer"
	FieldNameOperationFirst                 string = "gcp:operation_first"
	FieldNameOperationLast                  string = "gcp:operation_last"
	FieldNameLabels                         string = "gcp:labels"
	FieldNameCacheLookup                    string = "gcp:cache:lookup"
	FieldNameCacheHit                       string = "gcp:cache:hit"
	FieldNameCacheValidatedWithOriginServer string = "gcp:cache:validation_with_origin_header"
	FieldNameCacheFillBytes                 string = "gcp:cache:fill_bytes"
	FieldNameServerIp                       string = "gcp:server_ip"
	FieldNameLatency                        string = "gcp:latency"
	FieldNameTraceId                        string = "gcp:trace:id"
	FieldNameTraceSpanId                    string = "gcp:trace:span_id"
	FieldNameTraceEnabled                   string = "gcp:trace:enabled"
)

// gcp cloud logging severity
const (
	SeverityDefault   string = "DEFAULT"   // The log entry has no assigned severity level.
	SeverityDebug     string = "DEBUG"     // Debug or trace information.
	SeverityInfo      string = "INFO"      // Routine information, such as ongoing status or performance.
	SeverityNotice    string = "NOTICE"    // Normal but significant events, such as start up, shut down, or a configuration change.
	SeverityWarning   string = "WARNING"   // Warning events might cause problems.
	SeverityError     string = "ERROR"     // Error events are likely to cause problems.
	SeverityCritical  string = "CRITICAL"  // Critical events cause more severe problems or outages.
	SeverityAlert     string = "ALERT"     // A person must take an action immediately.
	SeverityEmergency string = "EMERGENCY" // One or more systems are unusable.
)
