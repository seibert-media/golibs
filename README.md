# //S/M Go Libs

[![Go Report Card](https://goreportcard.com/badge/github.com/seibert-media/golibs)](https://goreportcard.com/report/github.com/seibert-media/golibs)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Build Status](https://travis-ci.org/seibert-media/golibs.svg?branch=master)](https://travis-ci.org/seibert-media/golibs)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/f61779459d564fb59fc1013d27b36b1f)](https://www.codacy.com/app/seibert-media/golibs?utm_source=github.com&utm_medium=referral&utm_content=seibert-media/golibs&utm_campaign=Badge_Grade)
[![Codacy Badge](https://api.codacy.com/project/badge/Coverage/f61779459d564fb59fc1013d27b36b1f)](https://www.codacy.com/app/seibert-media/golibs?utm_source=github.com&utm_medium=referral&utm_content=seibert-media/golibs&utm_campaign=Badge_Coverage)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/seibert-media/golibs)

The repository contains various shared libs for use in //SEIBERT/MEDIA Golang projects.

## Libs

### Logging

This logging setup is using go.uber.org/zap.
Sentry is being added for production environments.

The signature of `log.New()` allows setting the Sentry DSN as well as a boolean for determining if we are in local mode.

```go
l := log.New(
    "sentryDSN",
    false,
)
```

Afterwards the logger can be used just like a default zap.Logger.

#### Modes

- `local=true`: In local mode, all logs will be printed in a human readable format.
- `local=false`: In non-local mode, all logs will be printed in [Stackdriver](https://cloud.google.com/logging/) format.
- `sentryDSN=""`: If no sentry dsn is set, there will be no error reporting.
- `sentryDSN="valid sentry dsn"`: If a valid sentry dsn is set, all logs >= Error will get reported to [Sentry](https://sentry.io).

#### Log Levels

This logging setup supports zap's dynamic log level.
By default, the logger will only print Info and above.

To change the level on the fly, use the provided helpers:

```go
logger, _ = log.New("", true)
logger.SetLevel(zap.DebugLevel)
logger.Debug("test")
```

For loggers stored in context, there is an additional helper:

```go
logger, _ := log.New("", true)
ctx := logger.To(context.Background())
log.SetLevel(ctx, zap.DebugLevel)
logger.Debug("test")
```

#### Sentry

To directly access Sentry the internal client is public.
All log keys will get sent to Sentry accordingly. Stacktraces will get attached to every Sentry Message.

#### Adding Fields to the Logger/Sentry

After initialization, the logger can be injected with fields which then get added to every log entry.
This is basic zap functionality, while being wrapped by our function to add those fields to Sentry as well.

Initializing a new logger and adding an app and a version field could look like this

```go
logger := log.New(*sentryDsn, *dbg).WithFields([]zapcore.Field{
    zap.String("app", "foobar"),
    zap.String("version", "0.1"),
})
defer logger.Sync()
```

#### Using with Context

Like the previous versions of this library, this one is primarily meant to be used in combination with `context.Context`.
This way we are able to pass logging down all our execution trees without ever making the logger a real dependency to call a function or initialize a struct.

Passing the previously defined logger into the root context is fairly simple and a helper function to do so is already available.

```go
ctx := log.WithLogger(context.Background(), logger)
```

This will return a new context containing the passed in logger based on a fresh context (`context.Background()`). Another context can be used instead
if available.

To then call the logger inside functions where said context gets passed into, simply use the available helper function for retrieving it.
Then use the logger as usual.

```go
log.From(ctx).Info("preparing")
log.From(ctx).Error("that did not work", zap.String("foo", "bar"), zap.Error(err))
// Sentry is available this way as well, always verify the client is not nil in case you disabled it
log.From(ctx).Sentry.SetEnvironment("dev")
```

Additionally there is a helper for adding new fields to the logger directly from context.

```go
ctx = log.WithFields(ctx, zap.String("newField", "value"))
```

#### [Experimental] Adding Sentry Release Info

Since the last version, it is supported to add Sentry release information to the logger.
This allows release tracking inside Sentry and therefor more insight into errors.

**INFO: For now release info has to be the first change made to the logger, right after setting fields. Otherwise the release info might not be persisted. Changing the release in a later stage is not possible.**

Example:

```go
logger := log.New("sentryDSN", false).WithRelease("some commit hash")
logger = logger.WithFields(zap.String("app", "example app"))
defer logger.Sync()
```

## Compatibility

This library requires at least Go 1.9+ and is currently tested against Go 1.9.x, 1.10.x and 1.11.x
For an up-to-date status on this check [.travis.yml](.travis.yml).

## Contributions

Pull Requests and Issue Reports are welcome.
