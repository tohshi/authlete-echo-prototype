package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/authlete/authlete-go/api"
	"github.com/authlete/authlete-go/conf"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var authleteApi api.AuthleteApi
var logger *zap.Logger

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func init() {
	l, _ := zap.NewProduction()
	logger = l

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	conf := new(conf.AuthleteEnvConfiguration)
	authleteApi = api.New(conf)
}

func main() {
	e := echo.New()

	e.Use(zapMiddleware(logger))
	e.Use(middleware.RequestID())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))))
	e.Use(middleware.Recover())

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("templates/*")),
	}

	e.GET(PING_ENDPOINT, func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "OK")
	})

	e.GET(AUTHORIZATION_ENDPOINT, authorizationHandler)
	e.POST(AUTHORIZATION_ENDPOINT, authorizationHandler)
	e.POST(TOKEN_ENDPOINT, tokenHandler)
	e.GET(LOGIN_ENDPOINT, loginPageHandler)
	e.POST(LOGIN_ENDPOINT, loginAttemptHandler)
	e.GET(CONSENT_ENDPOINT, consentPageHandler)
	e.POST(CONSENT_ENDPOINT, consentAttemptHandler)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func zapMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	})
}
