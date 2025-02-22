package mocks

import (
	"github.com/sing3demons/go-backend-clean-architecture/bootstrap"
	"github.com/stretchr/testify/mock"
)

type MockKafkaContext struct {
	mock.Mock
}

func (m *MockKafkaContext) Log(message string) {
	m.Called(message)
}

func (m *MockKafkaContext) Param(name string) string {
	args := m.Called(name)
	return args.String(0)
}

func (m *MockKafkaContext) Query(name string) string {
	args := m.Called(name)
	return "Query: " + args.String(0)
}

func (m *MockKafkaContext) ReadInput(data any) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockKafkaContext) Response(code int, data any) error {
	args := m.Called(code, data)
	return args.Error(0)
}

func (m *MockKafkaContext) SendMessage(topic string, payload any, opts ...bootstrap.OptionProducerMsg) (bootstrap.RecordMetadata, error) {
	args := m.Called(topic, payload, opts)
	return args.Get(0).(bootstrap.RecordMetadata), args.Error(1)
}

func NewMockKafkaContext() *MockKafkaContext {
	mockKafkaContext := new(MockKafkaContext)

	mockKafkaContext.On("Log", mock.Anything).Return()
	mockKafkaContext.On("Param", mock.Anything).Return("Param")
	mockKafkaContext.On("Query", mock.Anything).Return("Query")
	mockKafkaContext.On("ReadInput", mock.Anything).Return(nil)
	mockKafkaContext.On("Response", mock.Anything, mock.Anything).Return(nil)

	return mockKafkaContext
}
