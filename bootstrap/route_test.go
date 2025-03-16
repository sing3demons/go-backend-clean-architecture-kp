package bootstrap

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	handlerCalledErr = "handler was not called"
)

func TestHttpApplicationGet(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}
	app := newServer(cfg, log).(*httpApplication)

	handlerCalled := false
	app.Get("/test", func(ctx IContext) error {
		handlerCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, handlerCalled, handlerCalledErr)
}

func TestHttpApplicationPost(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}
	app := newServer(cfg, log).(*httpApplication)

	handlerCalled := false
	app.Post("/test", func(ctx IContext) error {
		handlerCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, handlerCalled, handlerCalledErr)
}

func TestHttpApplicationPut(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}
	app := newServer(cfg, log).(*httpApplication)

	handlerCalled := false
	app.Put("/test", func(ctx IContext) error {
		handlerCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodPut, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, handlerCalled, handlerCalledErr)
}

func TestHttpApplicationDelete(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}
	app := newServer(cfg, log).(*httpApplication)

	handlerCalled := false
	app.Delete("/test", func(ctx IContext) error {
		handlerCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodDelete, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, handlerCalled, handlerCalledErr)
}

func TestHttpApplicationPatch(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}
	app := newServer(cfg, log).(*httpApplication)

	handlerCalled := false
	app.Patch("/test", func(ctx IContext) error {
		handlerCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodPatch, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Response() should return status code 200")
	assert.True(t, handlerCalled, handlerCalledErr)
}

func TestHttpApplicationUse(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}
	app := newServer(cfg, log).(*httpApplication)

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

	assert.Equal(t, http.StatusOK, rec.Code, "Response() should return status code 200")
	assert.True(t, handlerCalled, handlerCalledErr)
}

// Register
func TestHttpApplicationRegister(t *testing.T) {
	log := NewZapLogger(zap.NewNop())

	cfg := &Config{
		AppConfig: AppConfig{
			Port: "8888",
		},
	}
	app := newServer(cfg, log).(*httpApplication)

	server := app.Register()
	if server.Addr != ":8888" {
		t.Errorf("expected server address :8888, got %s", server.Addr)
	}
}
