package bootstrap

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestEchoApplicationGet(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}

	app := newEchoServer(cfg, log).(*echoApplication)

	handlerCalled := false

	app.Get("/test", func(ctx IContext) error {
		handlerCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, printErr(http.StatusOK, rec.Code))
	assert.True(t, handlerCalled, handlerCalledErr)
}

func TestEchoApplicationPost(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}

	app := newEchoServer(cfg, log).(*echoApplication)

	handlerCalled := false

	app.Post("/test", func(ctx IContext) error {
		handlerCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, printErr(http.StatusOK, rec.Code))
	assert.True(t, handlerCalled, handlerCalledErr)
}

func TestEchoApplicationPut(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}

	app := newEchoServer(cfg, log).(*echoApplication)

	handlerCalled := false

	app.Put("/test", func(ctx IContext) error {
		handlerCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodPut, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, printErr(http.StatusOK, rec.Code))
	assert.True(t, handlerCalled, handlerCalledErr)
}

func TestEchoApplicationDelete(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}

	app := newEchoServer(cfg, log).(*echoApplication)

	handlerCalled := false

	app.Delete("/test", func(ctx IContext) error {
		handlerCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodDelete, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, printErr(http.StatusOK, rec.Code))
	assert.True(t, handlerCalled, handlerCalledErr)
}

func TestEchoApplicationPatch(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}

	app := newEchoServer(cfg, log).(*echoApplication)

	handlerCalled := false

	app.Patch("/test", func(ctx IContext) error {
		handlerCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodPatch, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, printErr(http.StatusOK, rec.Code))
	assert.True(t, handlerCalled, handlerCalledErr)
}

func TestEchoApplicationUse(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}

	app := newEchoServer(cfg, log).(*echoApplication)

	handlerCalled := false

	app.Use(func(next HandleFunc) HandleFunc {
		return func(ctx IContext) error {
			handlerCalled = true
			return next(ctx)
		}
	})

	app.Get("/test", func(ctx IContext) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, printErr(http.StatusOK, rec.Code))
	assert.True(t, handlerCalled, handlerCalledErr)
}

func TestEchoApplicationRegister(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}

	app := newEchoServer(cfg, log).(*echoApplication)

	app.Get("/test", func(ctx IContext) error {
		return nil
	})

	server := app.Register()

	assert.NotNil(t, server)
}
