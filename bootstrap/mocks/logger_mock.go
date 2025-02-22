package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
	name string
}

func NewMockLogger() *MockLogger {
	mockLogger := new(MockLogger)

	// Set expectations for methods that will be called in tests
	mockLogger.On("Println", mock.Anything).Return()
	mockLogger.On("Info", mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything).Return()
	mockLogger.On("Error", mock.Anything).Return()
	mockLogger.On("Fatal", mock.Anything).Return()
	mockLogger.On("DPanic", mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything).Return()
	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Fatalf", mock.Anything, mock.Anything).Return()
	mockLogger.On("DPanicf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Err", mock.Anything, mock.Anything).Return()

	return mockLogger
}

func (m *MockLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockLogger) Debug(args ...any) {
	m.Called(args...)
}

func (m *MockLogger) Debugf(template string, args ...any) {
	m.Called(template, args)
}

func (m *MockLogger) Info(args ...any) {
	m.Called(args...)
}

func (m *MockLogger) Infof(template string, args ...any) {
	m.Called(template, args)
}

func (m *MockLogger) Warn(args ...any) {
	m.Called(args...)
}

func (m *MockLogger) Warnf(template string, args ...any) {
	m.Called(template, args)
}

func (m *MockLogger) WarnMsg(msg string, err error) {
	m.Called(msg, err)
}

func (m *MockLogger) Error(args ...any) {
	m.Called(args...)
}

func (m *MockLogger) Errorf(template string, args ...any) {
	m.Called(template, args)
}

func (m *MockLogger) Err(msg string, err error) {
	m.Called(msg, err)
}

func (m *MockLogger) DPanic(args ...any) {
	m.Called(args...)
}

func (m *MockLogger) DPanicf(template string, args ...any) {
	m.Called(template, args)
}

func (m *MockLogger) Fatal(args ...any) {
	m.Called(args...)
}

func (m *MockLogger) Fatalf(template string, args ...any) {
	m.Called(template, args)
}

func (m *MockLogger) Printf(template string, args ...any) {
	m.Called(template, args)
}

func (m *MockLogger) WithName(name string) {
	m.name = name
	m.Called(name)
}

func (m *MockLogger) Println(v ...any) {
	m.Called(v...)
}
