package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/IBM/sarama"
)

type kafkaContext struct {
	topic    string
	headers  map[string]string
	body     string
	producer sarama.SyncProducer
	Logger   ILogger
	ctx      context.Context
}

type OptionProducerMsg struct {
	key       string
	headers   []map[string]string
	Timestamp time.Time
	Metadata  any
	Offset    int64
	Partition int32
}

func newConsumer(option *KafkaConfig) (sarama.ConsumerGroup, error) {
	if option.consumer != nil {
		return option.consumer, nil
	}
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V2_5_0_0     // Ensure Kafka version compatibility
	config.Consumer.Return.Errors = true // Capture errors from Kafka

	if option.Username != "" && option.Password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = option.Username
		config.Net.SASL.Password = option.Password
	}

	return sarama.NewConsumerGroup(option.Brokers, option.GroupID, config)
}

// NewConsumerContext creates a new Kafka context for consumer
func NewConsumerContext(topic, body string, producer sarama.SyncProducer, log ILogger) IContext {
	ctx := InitSession(context.Background(), log)
	return &kafkaContext{
		topic:    topic,
		body:     body,
		producer: producer,
		Logger:   log,
		ctx:      ctx,
	}
}

func (c *kafkaContext) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

func (ctx *kafkaContext) Log() ILogger {
	switch logger := ctx.Context().Value(key).(type) {
	case ILogger:
		return logger
	default:
		return ctx.Logger
	}

}

func (ctx kafkaContext) Param(name string) string {
	if name == "topic" {
		return ctx.topic
	}
	return ""
}

func (ctx *kafkaContext) Query(name string) string {
	return ""
}

func (ctx *kafkaContext) SetHeader(key, value string) {
	if ctx.headers == nil {
		ctx.headers = make(map[string]string)
	}
	ctx.headers[key] = value
}

func (ctx *kafkaContext) GetHeader(key string) string {
	if ctx.headers == nil {
		return ""
	}
	return ctx.headers[key]
}

func (ctx *kafkaContext) ReadInput(data any) error {
	const errMsgFormat = "%s, payload: %s"
	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Ptr, reflect.Interface:
		if val.Elem().Kind() == reflect.String {
			val.Elem().SetString(ctx.body)
			return nil
		}

		if err := json.Unmarshal([]byte(ctx.body), data); err != nil {
			return fmt.Errorf(errMsgFormat, err.Error(), ctx.body)
		}
		return nil
	case reflect.String:
		return fmt.Errorf("cannot assign to non-pointer string")
	default:
		err := json.Unmarshal([]byte(ctx.body), &data)
		if err != nil {
			return fmt.Errorf(errMsgFormat, err.Error(), ctx.body)
		}
		return nil
	}
}

func (ctx *kafkaContext) Response(code int, data any) error {
	return nil
}

func (ctx *kafkaContext) SendMessage(topic string, payload any, opts ...OptionProducerMsg) (RecordMetadata, error) {
	return producer(ctx.producer, topic, payload, opts...)
}
