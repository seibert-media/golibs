package log_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
)

func Test_NewDebug(t *testing.T) {
	logger, _ := log.New("", true, true)
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
	l, _ := log.New("", true, true)
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
	l, _ := log.New("", true, true)
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
	l, _ := log.New("", true, true)
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
	l, _ := log.New("", true, true)
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
	logger, _ := log.New("", false, true)
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
	_, err := log.New("^", true, true)
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

	logger, _ = log.New("", false, true)
	logger = logger.WithRelease("test")
	if logger.Sentry.Release() != "test" {
		t.Fatal("sentry release info not set, is:", logger.Sentry.Release())
	}
}
