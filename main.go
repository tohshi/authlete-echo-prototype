package main

import (
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/authlete/authlete-go/api"
	"github.com/authlete/authlete-go/conf"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func init() {
}

func main() {
	e := echo.New()

	if err := godotenv.Load(); err != nil {
		e.Logger.Fatal("Error loading .env file")
	}

	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))))
	e.Use(middleware.Recover())

	conf := conf.AuthleteEnvConfiguration{}
	api := AuthleteApi{api.New(&conf)}
	e.Use(api.authleteMiddleware)

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("templates/*")),
	}

	e.GET(PING_ENDPOINT, func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "OK")
	})

	e.GET(AUTHORIZATION_ENDPOINT, authorizationHandler, api.authleteMiddleware)
	e.POST(AUTHORIZATION_ENDPOINT, authorizationHandler)
	e.POST(TOKEN_ENDPOINT, tokenHandler, api.authleteMiddleware)
	e.GET(LOGIN_ENDPOINT, loginPageHandler)
	e.POST(LOGIN_ENDPOINT, loginAttemptHandler)
	e.GET(CONSENT_ENDPOINT, consentPageHandler)
	e.POST(CONSENT_ENDPOINT, consentAttemptHandler)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
