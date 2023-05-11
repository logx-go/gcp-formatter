package gcpformatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/logx-go/commons/pkg/commons"
	"github.com/logx-go/contract/pkg/logx"
	"github.com/logx-go/gcp-formatter/pkg/gcpformatter/model"
)

var _ logx.Formatter = (*GCPFormatter)(nil)

// New returns a new GCP Cloud Logging compatible gcp_formatter
func New() *GCPFormatter {
	return &GCPFormatter{
		logLevelToSeverityMap: map[int]string{
			logx.LogLevelDebug:   SeverityDebug,
			logx.LogLevelInfo:    SeverityInfo,
			logx.LogLevelNotice:  SeverityNotice,
			logx.LogLevelWarning: SeverityWarning,
			logx.LogLevelError:   SeverityError,
			logx.LogLevelFatal:   SeverityAlert,
			logx.LogLevelPanic:   SeverityEmergency,
		},
		logLevelDefault: logx.LogLevelInfo,
	}
}

type GCPFormatter struct {
	logLevelToSeverityMap map[int]string
	logLevelDefault       int
	projectID             string
}

func (j *GCPFormatter) clone() *GCPFormatter {
	return &GCPFormatter{
		logLevelToSeverityMap: j.logLevelToSeverityMap,
		logLevelDefault:       j.logLevelDefault,
		projectID:             j.projectID,
	}
}

func (j *GCPFormatter) WithLogLevelToSeverityMap(m map[int]string) *GCPFormatter {
	c := j.clone()
	c.logLevelToSeverityMap = m

	return c
}

func (j *GCPFormatter) WithLogLevelDefault(l int) *GCPFormatter {
	c := j.clone()
	c.logLevelDefault = l

	return c
}

func (j *GCPFormatter) WithProjectID(p string) *GCPFormatter {
	c := j.clone()
	c.projectID = p

	return c
}

func (j *GCPFormatter) Format(message *string, fields map[string]any) (messageF string, fieldsF map[string]any) {
	data := &model.LogEntry{
		Severity:       j.formatSeverity(fields),
		InsertId:       commons.GetFieldAsStringOrElse(FieldNameInsertId, fields, ""),
		Trace:          commons.GetFieldAsStringOrElse(FieldNameTraceId, fields, ""),
		TraceSampled:   commons.GetFieldAsBoolOrElse(FieldNameTraceEnabled, fields, false),
		SpanId:         commons.GetFieldAsStringOrElse(FieldNameTraceSpanId, fields, ""),
		TextPayload:    commons.GetAsStringOrElse(message, ""),
		Labels:         commons.GetFieldAsStringMapOrElse(FieldNameLabels, fields, nil),
		Timestamp:      commons.GetFieldAsTimeOrElse(logx.FieldNameTimestamp, fields, time.Now()),
		JsonPayload:    j.formatJsonPayload(fields),
		HttpRequest:    j.formatHttpRequest(fields),
		Operation:      j.formatOperation(fields),
		SourceLocation: j.formatSourceLocation(fields),
	}

	j.formatTracing(fields, data)

	enc, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
	}

	return string(enc), fields
}

func (j *GCPFormatter) formatTracing(fields map[string]any, data *model.LogEntry) {
	req := commons.GetFieldAsRequestPtrOrElse(logx.FieldNameHTTPRequest, fields, nil)
	if data.Trace == "" && j.projectID != "" && req != nil {
		data.Trace = fmt.Sprintf(`projects/%s/traces/%s`, j.projectID, j.extractTraceID(req))
		data.SpanId = j.extractSpanID(req)
		data.TraceSampled = j.extractTraceEnabled(req)
	}
}

func (j *GCPFormatter) formatSourceLocation(fields map[string]any) *model.SourceLocation {
	sourceLocation := &model.SourceLocation{
		File:     commons.GetFieldAsStringOrElse(logx.FieldNameCallerFile, fields, ""),
		Line:     commons.GetFieldAsStringOrElse(logx.FieldNameCallerLine, fields, ""),
		Function: commons.GetFieldAsStringOrElse(logx.FieldNameCallerFunc, fields, ""),
	}

	if sourceLocation.File == "" {
		return nil
	}

	return sourceLocation
}

func (j *GCPFormatter) formatJsonPayload(fields map[string]any) map[string]json.RawMessage {
	hasEntries := false
	jsonPayload := make(map[string]json.RawMessage)
	skip := []string{
		logx.FieldNameCallerFile,
		logx.FieldNameCallerLine,
		logx.FieldNameCallerFunc,
		logx.FieldNameLogLevel,
		logx.FieldNameMessage,
		logx.FieldNameTimestamp,
		logx.FieldNameHTTPRequest,
		logx.FieldNameHTTPResponse,
		FieldNameTraceId,
		FieldNameTraceEnabled,
		FieldNameTraceSpanId,
		FieldNameServerIp,
		FieldNameCacheLookup,
		FieldNameCacheHit,
		FieldNameCacheValidatedWithOriginServer,
		FieldNameCacheFillBytes,
		FieldNameLatency,
		FieldNameInsertId,
		FieldNameLabels,
		FieldNameOperationId,
		FieldNameOperationProducer,
		FieldNameOperationFirst,
		FieldNameOperationLast,
	}

	for name, value := range fields {
		if commons.Contains(skip, name) {
			continue
		}

		if raw, err := json.Marshal(value); err == nil {
			jsonPayload[name] = raw
			hasEntries = true
		}
	}

	if !hasEntries {
		return nil
	}

	return jsonPayload
}

func (j *GCPFormatter) formatOperation(fields map[string]any) *model.Operation {
	opId := commons.GetFieldAsStringOrElse(FieldNameOperationId, fields, "")
	opProd := commons.GetFieldAsStringOrElse(FieldNameOperationProducer, fields, "")
	opFirst := commons.GetFieldAsBoolOrElse(FieldNameOperationFirst, fields, false)
	opLast := commons.GetFieldAsBoolOrElse(FieldNameOperationLast, fields, false)

	if opId == "" && opProd == "" {
		return nil
	}

	return &model.Operation{
		Id:       opId,
		Producer: opProd,
		First:    opFirst,
		Last:     opLast,
	}
}

func (j *GCPFormatter) formatHttpRequest(fields map[string]any) *model.HttpRequest {
	req := commons.GetFieldAsRequestPtrOrElse(logx.FieldNameHTTPRequest, fields, nil)
	if req == nil {
		return nil
	}

	result := &model.HttpRequest{
		RequestMethod:                  req.Method,
		RequestUrl:                     req.RequestURI,
		RequestSize:                    j.calculateRequestSize(req),
		UserAgent:                      req.UserAgent(),
		RemoteIp:                       req.RemoteAddr,
		ServerIp:                       commons.GetFieldAsStringOrElse(FieldNameServerIp, fields, ""),
		Protocol:                       req.Proto,
		Referer:                        req.Referer(),
		CacheLookup:                    commons.GetFieldAsBoolOrElse(FieldNameCacheLookup, fields, false),
		CacheHit:                       commons.GetFieldAsBoolOrElse(FieldNameCacheHit, fields, false),
		CacheValidatedWithOriginServer: commons.GetFieldAsBoolOrElse(FieldNameCacheValidatedWithOriginServer, fields, false),
		CacheFillBytes:                 commons.GetFieldAsStringOrElse(FieldNameCacheValidatedWithOriginServer, fields, ""),
		Latency:                        commons.GetFieldAsStringOrElse(FieldNameLatency, fields, ""),
	}

	res := commons.GetFieldAsResponsePtrOrElse(logx.FieldNameHTTPResponse, fields, nil)
	if nil == res {
		return result
	}

	result.Status = res.StatusCode
	result.ResponseSize = j.calculateResponseSize(res)

	return result
}

func (j *GCPFormatter) calculateResponseSize(resp *http.Response) string {
	statusLineSize := int64(len(resp.Proto) + len(resp.Status) + 5) // space + status code (3 bytes) + CRLF (2 bytes)

	var headersSize int64
	for k, v := range resp.Header {
		for _, value := range v {
			headersSize += int64(len(k) + len(value) + 4) // ": " (2 bytes) + CRLF (2 bytes)
		}
	}
	headersSize += 2 // Final CRLF after headers (2 bytes)

	bodySize := int64(0)
	if resp.Body != nil {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, resp.Body)
		if err != nil {
			return ""
		}
		bodySize = int64(buf.Len())

		// Reset the response body to its original state
		resp.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	}

	return fmt.Sprintf(`%d`, statusLineSize+headersSize+bodySize)
}

func (j *GCPFormatter) calculateRequestSize(req *http.Request) string {
	requestLineSize := int64(len(req.Method) + len(req.URL.String()) + len(req.Proto) + 4) // 2 spaces + CRLF (2 bytes)

	var headersSize int64
	for k, v := range req.Header {
		for _, value := range v {
			headersSize += int64(len(k) + len(value) + 4) // ": " (2 bytes) + CRLF (2 bytes)
		}
	}
	headersSize += 2 // Final CRLF after headers (2 bytes)

	bodySize := int64(0)
	if req.Body != nil {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, req.Body)
		if err != nil {
			return ""
		}
		bodySize = int64(buf.Len())

		// Reset the request body to its original state
		req.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	}

	return fmt.Sprintf(`%d`, requestLineSize+headersSize+bodySize)
}

func (j *GCPFormatter) extractTraceID(req *http.Request) string {
	if req == nil {
		return ""
	}

	header := req.Header.Get("X-Cloud-Trace-Context")
	if header == "" {
		return ""
	}

	parts := strings.Split(header, "/")
	if len(parts) != 2 {
		return ""
	}

	traceID := parts[0]

	return traceID
}

func (j *GCPFormatter) extractSpanID(req *http.Request) string {
	if req == nil {
		return ""
	}

	header := req.Header.Get("X-Cloud-Trace-Context")
	if header == "" {
		return ""
	}

	parts := strings.Split(header, "/")
	if len(parts) != 2 {
		return ""
	}

	spanIDAndTraceTrue := strings.Split(parts[1], ";")
	if len(spanIDAndTraceTrue) != 2 {
		return ""
	}

	spanID := spanIDAndTraceTrue[0]
	return spanID
}

func (j *GCPFormatter) extractTraceEnabled(req *http.Request) bool {
	if req == nil {
		return false
	}

	header := req.Header.Get("X-Cloud-Trace-Context")
	if header == "" {
		return false
	}

	parts := strings.Split(header, "/")
	if len(parts) != 2 {
		return false
	}

	spanIDAndTraceTrue := strings.Split(parts[1], ";")
	if len(spanIDAndTraceTrue) != 2 {
		return false
	}

	traceTrue := spanIDAndTraceTrue[1] == "o=1"
	return traceTrue
}

func (j *GCPFormatter) formatSeverity(fields map[string]any) string {
	lvl := commons.GetFieldAsIntOrElse(logx.FieldNameLogLevel, fields, j.logLevelDefault)

	if s, ok := j.logLevelToSeverityMap[lvl]; ok {
		return s
	}

	return SeverityDefault
}
