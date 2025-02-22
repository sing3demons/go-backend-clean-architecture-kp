package bootstrap

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"go.uber.org/zap"
)

type FakeHttpContext struct {
	Res *httptest.ResponseRecorder
	Req *http.Request
	cfg *KafkaConfig
	log ILogger
}

func NewMockMuxContext(method string, url string, body any) *FakeHttpContext {
	var buf *bytes.Buffer
	if body != nil {
		jsonData, _ := json.Marshal(body)
		buf = bytes.NewBuffer(jsonData)
	} else {
		buf = &bytes.Buffer{}
	}

	req := httptest.NewRequest(method, url, buf)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rec := httptest.NewRecorder()

	// Create mock dependencies
	mockCfg := &KafkaConfig{}
	mockLog := NewZapLogger(zap.NewNop())

	return &FakeHttpContext{
		Res: rec,
		Req: req,
		cfg: mockCfg,
		log: mockLog,
	}
}

func (c *FakeHttpContext) Code() int {
	return c.Res.Code
}
func (c *FakeHttpContext) Body(data any) error {
	return json.NewDecoder(c.Res.Body).Decode(data)
}

func (c *FakeHttpContext) Context() context.Context {
	return c.Req.Context()
}

func (c *FakeHttpContext) SendMessage(topic string, message any, opts ...OptionProducerMsg) (RecordMetadata, error) {
	return producer(c.cfg.producer, topic, message, opts...)
}

func (c *FakeHttpContext) Log() ILogger {
	return c.log
}

func (c *FakeHttpContext) Query(name string) string {
	return c.Req.URL.Query().Get(name)
}

func (c *FakeHttpContext) Param(name string) string {
	v := c.Req.Context().Value(ContextKey(name))
	var value string
	switch v := v.(type) {
	case string:
		value = v
	}
	c.Req = c.Req.WithContext(context.WithValue(c.Req.Context(), ContextKey(name), nil))
	return value
}

func (c *FakeHttpContext) ReadInput(data any) error {
	return json.NewDecoder(c.Req.Body).Decode(data)
}

func (c *FakeHttpContext) Response(responseCode int, responseData any) error {
	c.Res.Header().Set("Content-type", "application/json; charset=UTF8")

	c.Res.WriteHeader(responseCode)

	err := json.NewEncoder(c.Res).Encode(responseData)
	return err
}
