# LogX - GCP Formatter

Google Cloud Logging compliant output formatter

## Install

```shell
go get -u github.com/logx-go/log-adapter
```

## Usage

```golang
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/logx-go/contract/pkg/logx"
	"github.com/logx-go/gcp-formatter/pkg/gcpformatter"
	"github.com/logx-go/log-adapter/pkg/logadapter"
)

func main() {
	// build/configure logger
	formatter := gcpformatter.New().
		// Optional: Will be used to compute the trace attribute (default: "")
		WithProjectID("your-gcp-project-id").
		// Optional: Default log level if non has been set (default: logx.LogLevelInfo)
		WithLogLevelDefault(logx.LogLevelDebug).
		// Optional: This is the default mapping from log level to severity
		WithLogLevelToSeverityMap(map[int]string{
			logx.LogLevelDebug:   gcpformatter.SeverityDebug,
			logx.LogLevelInfo:    gcpformatter.SeverityInfo,
			logx.LogLevelNotice:  gcpformatter.SeverityNotice,
			logx.LogLevelWarning: gcpformatter.SeverityWarning,
			logx.LogLevelError:   gcpformatter.SeverityError,
			logx.LogLevelFatal:   gcpformatter.SeverityAlert,
			logx.LogLevelPanic:   gcpformatter.SeverityEmergency,
		})

	// we are using the standard adapter from https://pkg.go.dev/log
	logger := logadapter.New(log.New(os.Stdout, "", 0)).
		WithFormatter(formatter)

	logSomething(logger)
}

func logSomething(logger logx.Logger) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req.Header.Set("X-Cloud-Trace-Context", "1c7886eaa2474d5da4da8c4f4bf6fdeb/1234567890;o=1")
	logger.Info("This is an error message",
		logx.FieldNameHttpRequest, req,
		"a random", "value",
		"number", 1981,
	)
}
```

This will output a single-line JSON to stdout:
```json
{
  "severity": "INFO",
  "traceSampled": true,
  "trace": "projects/your-gcp-project-id/traces/1c7886eaa2474d5da4da8c4f4bf6fdeb",
  "spanId": "1234567890",
  "textPayload": "This is an error message",
  "jsonPayload": {
    "a random": "value",
    "number": 1981
  },
  "httpRequest": {
    "requestMethod": "GET",
    "requestSize": "108",
    "protocol": "HTTP/1.1"
  },
  "sourceLocation": {
    "file": "/Users/mr/devel/logx/gcp-formatter/examples/basic.go",
    "line": "38",
    "function": "main.logSomething"
  },
  "timestamp": "2023-05-10T18:56:46.729533+02:00"
}
```
## Development

### Requirement
- Golang >=1.20
- golangci-lint (https://golangci-lint.run/)

### Tests

```shell
go test ./... -race
```

### Lint

```shell
golangci-lint run
```

## License

MIT License (see [LICENSE](LICENSE) file)

