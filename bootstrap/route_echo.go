package bootstrap

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type echoApplication struct {
	router      *echo.Echo
	middlewares []Middleware
	cfg         *Config
	log         ILogger
}

func newEchoServer(cfg *Config, log ILogger) IRouter {
	app := echo.New()

	return &echoApplication{
		router: app,
		cfg:    cfg,
		log:    log,
	}
}

func (app *echoApplication) Get(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.GET(path, func(c echo.Context) error {
		return preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newEchoContext(c, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *echoApplication) Post(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.POST(path, func(c echo.Context) error {
		return preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newEchoContext(c, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *echoApplication) Put(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.PUT(path, func(c echo.Context) error {
		return preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newEchoContext(c, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *echoApplication) Delete(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.DELETE(path, func(c echo.Context) error {
		return preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newEchoContext(c, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *echoApplication) Patch(path string, handler HandleFunc, middlewares ...Middleware) {
	app.router.PATCH(path, func(c echo.Context) error {
		return preHandle(handler, preMiddleware(app.middlewares, middlewares)...)(newEchoContext(c, &app.cfg.KafkaConfig, app.log))
	})
}

func (app *echoApplication) Use(middlewares ...Middleware) {
	app.middlewares = append(app.middlewares, middlewares...)
}

func (app *echoApplication) Register() *http.Server {
	return &http.Server{
		Addr:    ":" + app.cfg.AppConfig.Port,
		Handler: app.router,
	}
}

func (app *echoApplication) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.router.ServeHTTP(w, r)
	// app.router.NewContext(r, w)
}
