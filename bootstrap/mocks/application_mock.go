package mocks

import (
	"net/http"

	"github.com/sing3demons/go-backend-clean-architecture/bootstrap"
	"github.com/stretchr/testify/mock"
)

type MockApplication struct {
	mock.Mock
}

func (m *MockApplication) Get(path string, handler bootstrap.HandleFunc, middlewares ...bootstrap.Middleware) {
	m.Called(path, handler, middlewares)
}

func (m *MockApplication) Post(path string, handler bootstrap.HandleFunc, middlewares ...bootstrap.Middleware) {
	m.Called(path, handler, middlewares)
}

func (m *MockApplication) Use(middlewares ...bootstrap.Middleware) {
	m.Called(middlewares)
}

func (m *MockApplication) Start() {
	m.Called()
}

func (m *MockApplication) Consume(topic string, handler bootstrap.ServiceHandleFunc) {
	m.Called(topic, handler)
}

func (m *MockApplication) SendMessage(topic string, payload any, opts ...bootstrap.OptionProducerMsg) (bootstrap.RecordMetadata, error) {
	args := m.Called(topic, payload, opts)
	return args.Get(0).(bootstrap.RecordMetadata), args.Error(1)
}

func (m *MockApplication) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockApplication) Register() *http.Server {
	args := m.Called()
	return args.Get(0).(*http.Server)
}

func (m *MockApplication) Put(path string, handler bootstrap.HandleFunc, middlewares ...bootstrap.Middleware) {
	m.Called(path, handler, middlewares)
}

func (m *MockApplication) Delete(path string, handler bootstrap.HandleFunc, middlewares ...bootstrap.Middleware) {
	m.Called(path, handler, middlewares)
}

func (m *MockApplication) Patch(path string, handler bootstrap.HandleFunc, middlewares ...bootstrap.Middleware) {
	m.Called(path, handler, middlewares)
}

func NewMockApplication() *MockApplication {
	mockApplication := new(MockApplication)

	mockApplication.On("Get", mock.Anything, mock.Anything, mock.Anything).Return()
	mockApplication.On("Post", mock.Anything, mock.Anything, mock.Anything).Return()
	mockApplication.On("Use", mock.Anything).Return()
	mockApplication.On("Start").Return()
	mockApplication.On("Consume", mock.Anything, mock.Anything).Return()
	mockApplication.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).Return()

	return mockApplication
}
