package log_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
)

func Test_NewWithSentry(t *testing.T) {
	logger, err := log.New("http://test@localhost/test", true)
	if err != nil {
		t.Fatal("creating logger failed with:", err)
	}
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
}

func Test_NewWithoutSentry(t *testing.T) {
	logger, err := log.New("", true)
	if err != nil {
		t.Fatal("creating logger failed with:", err)
	}
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry != nil {
		t.Fatal("sentry is not nil")
	}
}

func Test_NewDebug(t *testing.T) {
	logger, err := log.New("http://test@localhost/test", true)
	if err != nil {
		t.Fatal("creating logger failed with:", err)
	}
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
	logger = logger.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

func Test_From(t *testing.T) {
	l, _ := log.New("", true)
	ctx := context.Background()

	ctx = log.WithLogger(ctx, l)
	if log.From(ctx).IsNop() {
		t.Fatal("logger should not be nop")
	}
	log.From(ctx).Debug("test", zap.String("test", "test"))

	ctx = context.Background()
	if !log.From(ctx).IsNop() {
		t.Fatal("logger should be nop")
	}
}

func Test_WithFields(t *testing.T) {
	l, _ := log.New("http://test@localhost/test", true)
	ctx := context.Background()

	ctx = log.WithLogger(ctx, l)
	if log.From(ctx).IsNop() {
		t.Fatal("logger should not be nop")
	}
	log.From(ctx).Debug("test", zap.String("test", "test"))
	log.From(ctx).Sentry.SetRelease("test")
	ctx = log.WithFields(ctx, zap.String("test-new-field", "test"))
	if log.From(ctx).Sentry.Release() != "test" {
		t.Fatal("sentry release info should stay consistent when adding fields")
	}
	log.From(ctx).Debug("test", zap.String("test", "test"))
}

func Test_WithFieldsOverwrite(t *testing.T) {
	l, _ := log.New("http://test@localhost/test", true)
	ctx := context.Background()

	ctx = log.WithLogger(ctx, l)
	if log.From(ctx).IsNop() {
		t.Fatal("logger should not be nop")
	}
	log.From(ctx).Debug("test", zap.String("test", "test"))

	log.WithFieldsOverwrite(ctx, zap.String("test-new-field", "test"))

	log.From(ctx).Debug("test", zap.String("test", "test"))
}

func Test_To(t *testing.T) {
	l, _ := log.New("http://test@localhost/test", true)
	ctx := context.Background()

	ctx = l.To(ctx)
	if log.From(ctx).IsNop() {
		t.Fatal("logger should not be nop")
	}
	log.From(ctx).Debug("test", zap.String("test", "test"))

	ctx = context.Background()
	if !log.From(ctx).IsNop() {
		t.Fatal("logger should be nop")
	}
}

func Test_NewNoDebug(t *testing.T) {
	logger, _ := log.New("http://test@localhost/test", true)
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
	logger = logger.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

func Test_NewNotLocal(t *testing.T) {
	logger, _ := log.New("http://test@localhost/test", false)
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
	logger = logger.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

func Test_NewInvalidSentryURL(t *testing.T) {
	_, err := log.New("^", true)
	if err == nil {
		t.Errorf("New() should have returned error")
	}
}

func Test_NewNop(t *testing.T) {
	logger := log.NewNop()
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger = logger.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

func Test_SetRelease(t *testing.T) {
	logger := log.NewNop()
	logger = logger.WithRelease("test")
	if logger.Sentry.Release() != "" {
		t.Fatal("noop logger shouldn't have release info", logger.Sentry.Release())
	}

	logger, _ = log.New("http://test@localhost/test", true)
	logger = logger.WithRelease("test")
	if logger.Sentry.Release() != "test" {
		t.Fatal("sentry release info not set, is:", logger.Sentry.Release())
	}
}

func Test_SetLevel(t *testing.T) {
	out := make(chan string)

	capture := &stdCapture{}
	capture.capture(out)
	logger, _ := log.New("", true)
	logger.Debug("test")
	capture.finish()
	msg := <-out
	if len(msg) > 0 {
		t.Fatal("logger should not print message to stdout, got:", msg)
	}

	capture = &stdCapture{}
	capture.capture(out)
	logger, _ = log.New("", true)
	logger.SetLevel(zap.DebugLevel)
	logger.Debug("test2")
	capture.finish()
	msg = <-out
	if len(msg) < 1 {
		t.Fatal("logger should print message to stdout, got:", msg)
	}
	if !strings.Contains(msg, "DEBUG") {
		t.Fatal("message should be DEBUG, got:", msg)
	}
}

func Test_CtxSetLevel(t *testing.T) {
	out := make(chan string)

	capture := &stdCapture{}
	capture.capture(out)
	logger, _ := log.New("", true)
	ctx := logger.To(context.Background())
	log.SetLevel(ctx, zap.DebugLevel)
	logger.Debug("test")
	capture.finish()
	msg := <-out
	if len(msg) < 1 {
		t.Fatal("logger should print message to stdout, got:", msg)
	}
	if !strings.Contains(msg, "DEBUG") {
		t.Fatal("message should be DEBUG, got:", msg)
	}
}

type stdCapture struct {
	stdout *os.File
	r, w   *os.File
	c      chan string
}

func (s *stdCapture) capture(to chan string) {
	s.stdout = os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(fmt.Sprint("creating pipe failed with:", err))
	}
	os.Stdout = w
	s.r, s.w = r, w

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		to <- buf.String()
	}()
}

func (s *stdCapture) finish() {
	s.w.Close()
	os.Stdout = s.stdout
}

