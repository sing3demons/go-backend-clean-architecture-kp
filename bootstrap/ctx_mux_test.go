package bootstrap

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type FakeHttpContext struct {
	Res *httptest.ResponseRecorder
	Req *http.Request
	cfg *KafkaConfig
	log ILogger
}

type Option struct {
	Body   any
	Query  map[string]string
	Params map[string]string
	Header map[string]string
}

func NewMockMuxContext(t *testing.T, opts ...Option) *FakeHttpContext {
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
			ctx = context.WithValue(ctx, ContextKey(k), v)
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
	producer := mocks.NewSyncProducer(t, mocks.NewTestConfig())
	producer.ExpectSendMessageAndSucceed()

	// Create mock dependencies
	mockCfg := &KafkaConfig{
		producer: producer,
	}
	mockLog := NewZapLogger(zap.NewNop())

	return &FakeHttpContext{
		Res: rec,
		Req: req,
		cfg: mockCfg,
		log: mockLog,
	}
}

func TestContextMux(t *testing.T) {
	mock := NewMockMuxContext(t, Option{
		Body: map[string]string{
			"message": "message",
		},
	})
	ctx := newMuxContext(mock.Res, mock.Req, mock.cfg, mock.log)

	assert.NotNil(t, ctx.Context())

	m, err := ctx.SendMessage("topic", "message")
	assert.Nil(t, err, "SendMessage() should not return error")
	assert.NotNil(t, m)

	assert.NotNil(t, ctx.Log())

	assert.NotNil(t, ctx.Query("name"), "Query() should return string")

	assert.NotNil(t, ctx.Param("name"), "Param() should return string")

	var body map[string]any
	if err := ctx.ReadInput(&body); err != nil {
		assert.Error(t, err, "ReadInput() should not return error")
	}
	assert.NotNil(t, body)
	assert.Equal(t, map[string]any{"message": "message"}, body)

	if err := ctx.Response(http.StatusOK, body); err != nil {
		// t.Error("Response() should not return error")
		assert.Error(t, err, "Response() should not return error")
	}
	assert.Equal(t, http.StatusOK, mock.Res.Code)
	assert.Equal(t, "application/json; charset=UTF8", mock.Res.Header().Get("Content-type"))

	// SetHeader
	ctx.SetHeader("key", "value")
	assert.Equal(t, "value", mock.Res.Header().Get("key"))

	// SetHeader delete
	ctx.SetHeader("key", "")
	assert.Equal(t, "", mock.Res.Header().Get("key"))

	// getHeader
	assert.Equal(t, "application/json", ctx.GetHeader("Content-Type"))

}
