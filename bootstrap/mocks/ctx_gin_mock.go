package mocks

import (
	"github.com/sing3demons/go-backend-clean-architecture/bootstrap"
	"github.com/stretchr/testify/mock"
)

type MockContext struct {
	mock.Mock
}

func (m *MockContext) SendMessage(topic string, message any, opts ...bootstrap.OptionProducerMsg) (bootstrap.RecordMetadata, error) {
	args := m.Called(topic, message, opts)
	return args.Get(0).(bootstrap.RecordMetadata), args.Error(1)
}

func (m *MockContext) Log(message string) {
	m.Called(message)
}

func (m *MockContext) Query(name string) string {
	args := m.Called(name)
	return args.String(0)
}

func (m *MockContext) Param(name string) string {
	args := m.Called("param", name)
	return args.String(1)
}

func (m *MockContext) ReadInput(data any) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockContext) Response(responseCode int, responseData any) error {
	args := m.Called(responseCode, responseData)
	return args.Error(0)
}

func NewMockContext() *MockContext {
	mockContext := new(MockContext)

	mockContext.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).Return()
	mockContext.On("Log", mock.Anything).Return()
	mockContext.On("Query", mock.Anything).Return("Query")
	mockContext.On("Param", mock.Anything).Return("Param")
	mockContext.On("ReadInput", mock.Anything).Return(nil)
	mockContext.On("Response", mock.Anything, mock.Anything).Return(nil)

	return mockContext
}
