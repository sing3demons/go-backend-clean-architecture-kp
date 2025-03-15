package bootstrap

import (
	"testing"

	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/assert"
)

func TestProducer(t *testing.T) {
	mockProducer := mocks.NewSyncProducer(t, nil) // Use mocks.NewSyncProducer

	topic := "test-topic"
	payload := map[string]string{"message": "hello world"}

	// data, _ := json.Marshal(payload)

	mockProducer.ExpectSendMessageAndSucceed() // Expect a successful send

	recordMetadata, err := producer(mockProducer, topic, payload)

	assert.NoError(t, err)
	assert.Equal(t, topic, recordMetadata.TopicName)
	assert.GreaterOrEqual(t, recordMetadata.Partition, int32(0))
	assert.GreaterOrEqual(t, recordMetadata.Offset, int64(0))
}
