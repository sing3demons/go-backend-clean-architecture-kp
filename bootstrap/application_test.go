package bootstrap

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockRouter struct {
	mock.Mock
}

func (m *MockRouter) Get(path string, handler func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) {
	m.Called(path, handler, middlewares)
}

func (m *MockRouter) Post(path string, handler func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) {
	m.Called(path, handler, middlewares)
}

func (m *MockRouter) Use(middlewares ...func(http.Handler) http.Handler) {
	m.Called(middlewares)
}

func (m *MockRouter) Register() *http.Server {
	args := m.Called()
	return args.Get(0).(*http.Server)
}
func TestNewApplication(t *testing.T) {
	mockRouter := new(MockRouter)

	logger := NewZapLogger(zap.NewNop())

	mockRouter.On("Register").Return(&http.Server{})

	config := &Config{
		AppConfig: AppConfig{
			Port: "3000",
		},
	}

	app := NewApplication(config, logger)

	assert.NotNil(t, app, "Application should not be nil")
}
 