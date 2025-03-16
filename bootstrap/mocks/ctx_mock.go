package mocks

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/sing3demons/go-backend-clean-architecture/bootstrap"
	"go.uber.org/zap"
)

type FakeHttpContext struct {
	Res *httptest.ResponseRecorder
	Req *http.Request
	cfg *bootstrap.KafkaConfig
	log bootstrap.ILogger
}

type Option struct {
	Body   any
	Query  map[string]string
	Params map[string]string
	Header map[string]string
}

func NewMockMuxContext(opts ...Option) *FakeHttpContext {
	opt := &Option{}
	if len(opts) > 0 {
		opt = &opts[0]
	}

	// Create request
	var buf *bytes.Buffer
	if opt.Body != nil {
		jsonData, _ := json.Marshal(opt.Body)
		buf = bytes.NewBuffer(jsonData)
	} else {
		buf = &bytes.Buffer{}
	}

	req := httptest.NewRequest(http.MethodOptions, "/api", buf)

	if opt.Query != nil {
		u := url.Values{}
		for k, v := range opt.Query {
			u.Set(k, v)
		}
		req.URL.RawQuery = u.Encode()
	}

	if opt.Params != nil {
		ctx := req.Context()
		for k, v := range opt.Params {
			ctx = context.WithValue(ctx, bootstrap.ContextKey(k), v)
		}
		req = req.WithContext(ctx)
	}

	if opt.Header != nil {
		for k, v := range opt.Header {
			req.Header.Set(k, v)
		}
	} else {
		req.Header.Set("Content-Type", "application/json")
	}

	// Create response recorder
	rec := httptest.NewRecorder()

	// Create mock dependencies
	mockCfg := &bootstrap.KafkaConfig{}
	mockLog := bootstrap.NewZapLogger(zap.NewNop())

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

func (c *FakeHttpContext) SendMessage(topic string, message any, opts ...bootstrap.OptionProducerMsg) (bootstrap.RecordMetadata, error) {
	return bootstrap.RecordMetadata{
		TopicName: topic,
		Offset:    0,
		Partition: 0,
	}, nil
}

func (c *FakeHttpContext) Log() bootstrap.ILogger {
	return c.log
}

func (c *FakeHttpContext) Query(name string) string {
	return c.Req.URL.Query().Get(name)
}

func (c *FakeHttpContext) Param(name string) string {
	v := c.Req.Context().Value(bootstrap.ContextKey(name))
	var value string
	switch v := v.(type) {
	case string:
		value = v
	}
	c.Req = c.Req.WithContext(context.WithValue(c.Req.Context(), bootstrap.ContextKey(name), nil))
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

func (c *FakeHttpContext) SetHeader(key, value string) {
	if value == "" {
		c.Res.Header().Del(key)
		return
	}
	c.Res.Header().Set(key, value)
}

func (c *FakeHttpContext) GetHeader(key string) string {
	return c.Req.Header.Get(key)
}
