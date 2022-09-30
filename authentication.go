package main

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func loginPageHandler(ctx echo.Context) error {
	sess, err := session.Get(AUTHORIZATION_SESSION, ctx)
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}
	if sess.Values[TICKET] == nil {
		return echo.ErrBadRequest
	}

	return ctx.Render(http.StatusOK, LOGIN_TEMPLATE, map[string]string{"action": LOGIN_ENDPOINT})
}

func loginAttemptHandler(ctx echo.Context) error {
	authzSess, err := session.Get(AUTHORIZATION_SESSION, ctx)
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}
	authnSess, err := session.Get(AUTHENTICATION_SESSION, ctx)
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}

	userId := ctx.FormValue("user_id")
	password := ctx.FormValue("password")

	user := findUser(userId)
	if user.Id == userId && user.Password == password {
		authnSess.Values[USER_ID] = user.Id
		authnSess.Options.MaxAge = (60 * 60 * 24) * 30 // One Month
		if err := authnSess.Save(ctx.Request(), ctx.Response()); err != nil {
			ctx.Logger().Error(err)
		}
		if isUserConsented(userId, authzSess.Values[CLIENT_ID].(string)) {
			return authorizationIssueCaller(ctx)
		} else {
			return ctx.Redirect(http.StatusFound, CONSENT_ENDPOINT)
		}
	}

	return ctx.Render(http.StatusOK, LOGIN_TEMPLATE, map[string]string{"action": LOGIN_ENDPOINT, "message": "Login Failed"})
}
