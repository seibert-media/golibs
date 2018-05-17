# //S/M Go Libs
[![Go Report Card](https://goreportcard.com/badge/github.com/seibert-media/golibs)](https://goreportcard.com/report/github.com/seibert-media/golibs)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Build Status](https://travis-ci.org/seibert-media/golibs.svg?branch=master)](https://travis-ci.org/seibert-media/golibs)

The repository containing various shared libs for //SEIBERT/MEDIA projects.

## Libs

### Logging
Our logging setup using go.uber.org/zap.
Sentry and Jaeger are being added for production environments.

```go
l := log.New(
    "name",
    "sentryDSN",
    false,
)
defer l.Close()
```

Afterwards the logger can be used just like a default zap.Logger.
When the log level is Error or worse, a sentry message is being sent containing all string and int tags.
If you provide a zap.Error tag, the related stacktrace will also be attached.

Additionally there is a tracer(opentracing/jaeger) available in the logger which should be closed before exiting main.

## Contributions

Pull Requests and Issue Reports are welcome.
