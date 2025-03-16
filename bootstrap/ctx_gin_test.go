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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type FakeGinContext struct {
	Res *httptest.ResponseRecorder
	Req *http.Request
	cfg *KafkaConfig
	Log ILogger
	Ctx *gin.Context
}

func NewGinMuxContext(t *testing.T, opts ...Option) *FakeGinContext {
	gin.SetMode(gin.TestMode)
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

	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	return &FakeGinContext{
		Res: rec,
		Req: req,
		cfg: mockCfg,
		Log: mockLog,
		Ctx: c,
	}
}

func TestGinContext(t *testing.T) {

	mock := NewGinMuxContext(t, Option{
		Body: map[string]string{
			"message": "message",
		},
		Query: map[string]string{
			"name": "x",
		},
		Params: map[string]string{
			"id": "1",
		},
	})

	ctx := newGinContext(mock.Ctx, mock.cfg, mock.Log)

	assert.NotNil(t, ctx.Context())
	assert.NotNil(t, ctx.Log(), "log")

	assert.Equal(t, "x", ctx.Query("name"), "name")
	assert.Equal(t, "", ctx.Param("id"), "message")

	var data map[string]string
	err := ctx.ReadInput(&data)
	assert.Nil(t, err, "ReadInput() should not return error")
	assert.Equal(t, "message", data["message"], "message")

	// Test SendMessage
	m, err := ctx.SendMessage("topic", "message")
	assert.Nil(t, err, "SendMessage() should not return error")
	assert.NotNil(t, m, "SendMessage() should return RecordMetadata")

	err = ctx.Response(http.StatusOK, data)
	assert.Nil(t, err, "Response() should not return error")
	assert.Equal(t, http.StatusOK, mock.Res.Code, "Response() should return status code 200")

	// Test SetHeader
	ctx.SetHeader("key", "value")
	assert.Equal(t, "value", mock.Res.Header().Get("key"), "SetHeader() should set header key")

	// Test GetHeader
	assert.Equal(t, "application/json", ctx.GetHeader("Content-Type"), "GetHeader() should return header value")
}
