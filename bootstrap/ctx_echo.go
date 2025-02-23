package bootstrap

import (
	"context"

	"github.com/labstack/echo/v4"
)

type EchoContext struct {
	ctx echo.Context
	cfg *KafkaConfig
	log ILogger
}

func newEchoContext(c echo.Context, cfg *KafkaConfig, log ILogger) IContext {
	ctx := InitSession(c.Request().Context(), log)
	c.Request().WithContext(ctx)
	return &EchoContext{ctx: c, cfg: cfg, log: log}
}

func (c *EchoContext) Context() context.Context {
	return c.ctx.Request().Context()
}

func (c *EchoContext) SendMessage(topic string, message any, opts ...OptionProducerMsg) (RecordMetadata, error) {
	return producer(c.cfg.producer, topic, message, opts...)
}

func (c *EchoContext) Log() ILogger {
	return c.log
}

func (c *EchoContext) Query(name string) string {
	return c.ctx.QueryParam(name)
}

func (c *EchoContext) Param(name string) string {
	return c.ctx.Param(name)
}

func (c *EchoContext) ReadInput(data any) error {
	return c.ctx.Bind(data)
}

func (c *EchoContext) Response(code int, data any) error {
	return c.ctx.JSON(code, data)
}

func (c *EchoContext) SetHeader(key, value string) {
	c.ctx.Response().Header().Set(key, value)
}

func (c *EchoContext) GetHeader(key string) string {
	return c.ctx.Request().Header.Get(key)
}
