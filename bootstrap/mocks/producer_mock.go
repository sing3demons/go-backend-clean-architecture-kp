package mocks

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/mock"
)

type MockSyncProducer struct {
	mock.Mock
}

func (m *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	args := m.Called(msg)
	return args.Get(0).(int32), args.Get(1).(int64), args.Error(2)
}

func (m *MockSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	args := m.Called(msgs)
	return args.Error(0)
}

func (m *MockSyncProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func NewMockSyncProducer() *MockSyncProducer {
	return &MockSyncProducer{}
}
