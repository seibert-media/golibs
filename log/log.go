package log

import (
	"context"
	"fmt"
	"os"

	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// Context implements context.Context while adding our own logging and tracing functionality
type Context struct {
	context.Context
	*zap.Logger
	Sentry *raven.Client

	dsn   string
	debug bool
}

// New Context with included logger and sentry instances
func New(ctx context.Context, dsn string, debug bool) *Context {
	sentry, err := raven.New(dsn)
	if err != nil {
		panic(err)
	}

	logger := buildLogger(sentry, debug)

	return &Context{
		Context: ctx,
		Logger:  logger,
		Sentry:  sentry,

		dsn:   dsn,
		debug: debug,
	}
}

// NewNop returns Context with empty logging and tracing
func NewNop(ctx context.Context) *Context {
	sentry, _ := raven.New("")
	logger := zap.NewNop()

	log := &Context{
		Context: ctx,
		Logger:  logger,
		Sentry:  sentry,
	}

	return log
}

// WithFields wrapper around zap.With
func (c *Context) WithFields(fields ...zapcore.Field) *Context {
	l := New(c.Context, c.dsn, c.debug)
	l.Logger = l.Logger.With(fields...)
	return l
}

// NewSentryEncoder with dsn
func NewSentryEncoder(client *raven.Client) zapcore.Encoder {
	return newSentryEncoder(client)
}

func newSentryEncoder(client *raven.Client) *sentryEncoder {
	enc := &sentryEncoder{}
	enc.Sentry = client
	return enc
}

type sentryEncoder struct {
	zapcore.ObjectEncoder
	dsn    string
	Sentry *raven.Client
}

// Clone .
func (s *sentryEncoder) Clone() zapcore.Encoder {
	return newSentryEncoder(s.Sentry)
}

// EncodeEntry .
func (s *sentryEncoder) EncodeEntry(e zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf := buffer.NewPool().Get()
	if e.Level == zapcore.ErrorLevel {
		tags := make(map[string]string)
		var err error
		for _, f := range fields {
			var tag string
			switch f.Type {
			case zapcore.StringType:
				tag = f.String
			case zapcore.Int16Type, zapcore.Int32Type, zapcore.Int64Type:
				tag = fmt.Sprintf("%v", f.Integer)
			case zapcore.ErrorType:
				err = f.Interface.(error)
			}
			tags[f.Key] = tag

		}
		if err == nil {
			s.Sentry.CaptureMessage(e.Message, tags)
			return buf, nil
		}
		s.Sentry.CaptureError(errors.Wrap(err, e.Message), tags)
	}
	return buf, nil
}

func (s *sentryEncoder) AddString(key, val string) {
	tags := s.Sentry.Tags
	if tags == nil {
		tags = make(map[string]string)
	}
	tags[key] = val
	s.Sentry.SetTagsContext(tags)
}

func (s *sentryEncoder) AddInt64(key string, val int64) {
	tags := s.Sentry.Tags
	if tags == nil {
		tags = make(map[string]string)
	}
	tags[key] = fmt.Sprint(val)
	s.Sentry.SetTagsContext(tags)
}

// buildLogger
func buildLogger(sentry *raven.Client, debug bool) *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel && lvl < zapcore.InfoLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	consoleConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)
	sentryEncoder := NewSentryEncoder(sentry)
	var core zapcore.Core
	if debug {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
			zapcore.NewCore(consoleEncoder, consoleDebugging, debugPriority),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
			zapcore.NewCore(sentryEncoder, consoleErrors, highPriority),
		)
	}

	logger := zap.New(core)
	if debug {
		logger = logger.WithOptions(
			zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel),
		)
	} else {
		logger = logger.WithOptions(
			zap.AddStacktrace(zap.FatalLevel),
		)
	}
	return logger
}
