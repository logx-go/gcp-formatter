package gcpformatter

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/logx-go/contract/pkg/logx"
	"github.com/stretchr/testify/assert"
)

func TestGCPFormatter_Format(t *testing.T) {
	f := New().WithProjectID("test")

	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req.Header.Set("X-Cloud-Trace-Context", "1c7886eaa2474d5da4da8c4f4bf6fdeb/1234567890;o=1")

	ts := time.Now()
	assert.Equal(t, fmt.Sprintf(
		`{"foo":"bar","httpRequest":{"requestMethod":"GET","requestUrl":"https://example.com","requestSize":"108","protocol":"HTTP/1.1"},"message":"test","severity":"INFO","sourceLocation":{"file":"file","line":"123","function":"func"},"spanId":"1234567890","timestamp":"%s","trace":"projects/test/traces/1c7886eaa2474d5da4da8c4f4bf6fdeb","traceSampled":true}`,
		ts.Format(time.RFC3339Nano),
	),
		f.Format("test", map[string]any{
			"foo":                     "bar",
			logx.FieldNameHTTPRequest: req,
			logx.FieldNameCallerFile:  "file",
			logx.FieldNameCallerFunc:  "func",
			logx.FieldNameCallerLine:  "123",
			logx.FieldNameTimestamp:   ts,
		}))
}
