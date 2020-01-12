crzerolog
===
[![godoc](https://godoc.org/github.com/yfuruyama/crzerolog?status.svg)](https://godoc.org/github.com/yfuruyama/crzerolog) [![CircleCI](https://circleci.com/gh/yfuruyama/crzerolog.svg?style=svg)](https://circleci.com/gh/yfuruyama/crzerolog)

A zerolog-based logging library for Cloud Run.

![request log](img/request_log.png)

## Features
- Auto format for Stackdriver fields such as time, severity, trace, sourceLocation
- Group application logs with the request log
- Supports all of [rs/zerolog](https://github.com/rs/zerolog) APIs for structured logging

## Installation

```
go get -u github.com/yfuruyama/crzerolog
```

## Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/yfuruyama/crzerolog"
)

func main() {
	rootLogger := zerolog.New(os.Stdout)
	middleware := crzerolog.InjectLogger(&rootLogger)

	http.Handle("/", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := hlog.FromRequest(r)

		logger.Info().Msg("Hi")
		logger.Warn().Str("foo", "bar").Msg("This is")
		logger.Error().Int("num", 123).Msg("Structured Log")

		fmt.Fprintf(w, "Hello\n")
	})))

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	log.Printf("Server listening on port %q", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

After running above code on your Cloud Run service, you can find following logs in Stackdriver Logging.

### Request Log
![request log](img/request_log.png)

### Container Logs
As well as request log, container log is also printed out. Note that some Stackdriver fields are automatically set, such as `severity`, `sourceLocation`, `spanId`, `timestamp`, `trace`.

![container log 1](img/container_log_01.png)

If you add additional fields to the log with zerolog APIs, you can find those fields in `jsonPayload`.

![container log 2](img/container_log_02.png)

![container log 3](img/container_log_03.png)

## Level mapping
This library automatically maps [zerolog level](https://godoc.org/github.com/rs/zerolog#Level) to [Stackdriver severity](https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity).

Mapping is as follows.

| zerolog level | Stackdriver severity |
| --- | --- |
| NoLevel | DEFAULT |
| TraceLevel | DEFAULT |
| DebugLevel | DEBUG |
| InfoLevel | INFO |
| WarnLevel | WARNING |
| ErrorLevel | ERROR |
| FatalLevel | CRITICAL |
| PanicLevel | ALERT |

## Supported Platform
- Cloud Run (fully managed)
- Google App Engine (2nd-Generation)

## License
[Apache 2.0](LICENSE).

## Disclaimer
This is not an official Google product.
