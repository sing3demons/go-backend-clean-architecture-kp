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
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type FakeEchoContext struct {
	Res *httptest.ResponseRecorder
	Req *http.Request
	cfg *KafkaConfig
	Log ILogger
	Ctx echo.Context
}

func NewEchoMuxContext(t *testing.T, opts ...Option) *FakeEchoContext {
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
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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

	c := echo.New().NewContext(req, rec)

	return &FakeEchoContext{
		Res: rec,
		Req: req,
		cfg: mockCfg,
		Log: mockLog,
		Ctx: c,
	}
}

func TestEchoContext(t *testing.T) {
	mock := NewEchoMuxContext(t, Option{
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
	ctx := newEchoContext(mock.Ctx, mock.cfg, mock.Log)

	assert.NotNil(t, ctx)
	assert.NotNil(t, ctx.Context())
	assert.NotNil(t, ctx.Log(), "Log() should return ILogger")
	assert.NotNil(t, ctx.Query("name"), "Query() should return string")
	assert.NotNil(t, ctx.Param("id"), "Param() should return string")

	var body map[string]any
	if err := ctx.ReadInput(&body); err != nil {
		assert.Error(t, err, "ReadInput() should not return error")
	}

	err := ctx.Response(http.StatusOK, body)
	assert.Nil(t, err, "Response() should not return error")

	m, err := ctx.SendMessage("topic", "message")
	assert.Nil(t, err, "SendMessage() should not return error")
	assert.NotNil(t, m)

	ctx.SetHeader("key", "value")
	assert.Equal(t, "value", mock.Res.Header().Get("key"))

	assert.Equal(t, http.StatusOK, mock.Res.Code)
	assert.Equal(t, echo.MIMEApplicationJSON, mock.Res.Header().Get(echo.HeaderContentType))

	ctx.SetHeader("key", "")
	assert.Equal(t, "", mock.Res.Header().Get("key"))

	assert.Equal(t, echo.MIMEApplicationJSON, ctx.GetHeader(echo.HeaderContentType))

}
