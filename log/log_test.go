package log_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
)

func Test_NewDebug(t *testing.T) {
	ctx := log.New(context.Background(), "", true)
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Logger == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	ctx.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
	ctx = ctx.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Logger == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	ctx.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

func Test_NewNoDebug(t *testing.T) {
	ctx := log.New(context.Background(), "", false)
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Logger == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	ctx.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
	ctx = ctx.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Logger == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	ctx.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

func Test_NewInvalidSentryURL(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("New() should have panicked")
			}
		}()
		log.New(context.Background(), "^", true)
	}()
}

func Test_NewNop(t *testing.T) {
	ctx := log.NewNop(context.Background())
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Logger == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	ctx = ctx.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Logger == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	ctx.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

type TestCtxKey string

func Test_ContextWorks(t *testing.T) {
	ctx := log.New(context.Background(), "", true)
	ctx.WithValue(TestCtxKey("test"), "test")
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	if ctx.Value(TestCtxKey("test")) != "test" {
		t.Fatal("ctx should contain text")
	}
	ctx, cancel := ctx.WithCancel()
	if ctx.Err() != nil {
		t.Fatal("context should not have error")
	}
	cancel()
	if ctx.Err() != context.Canceled {
		t.Fatal("context should be closed")
	}
	ctx = log.New(context.Background(), "", true)
	ctx.WithDeadline(time.Now().Add(1 * time.Millisecond))
	select {
	case <-time.After(2 * time.Millisecond):
		t.Fatal("context should be closed after deadline")
	case <-ctx.Done():
		break
	}
	ctx = log.New(context.Background(), "", true)
	ctx.WithTimeout(1 * time.Millisecond)
	select {
	case <-time.After(2 * time.Millisecond):
		t.Fatal("context should be closed after deadline")
	case <-ctx.Done():
		break
	}
}
