package main

import (
	"encoding/json"
	"net/http"

	"github.com/authlete/authlete-go/dto"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func isUserConsented(userId string, clientId string) bool {
	return userId == "user2"
}

func consentPageHandler(ctx echo.Context) error {
	sess, err := session.Get(AUTHORIZATION_SESSION, ctx)
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}
	if sess.Values[TICKET] == nil {
		return echo.ErrBadRequest
	}

	scopes := []dto.Scope{}
	json.Unmarshal(sess.Values[SCOPES].([]byte), &scopes)

	return ctx.Render(http.StatusOK, CONSENT_TEMPLATE, map[string]any{"action": CONSENT_ENDPOINT, "scopes": scopes})
}

func consentAttemptHandler(ctx echo.Context) error {
	var consent bool
	switch ctx.FormValue("consent") {
	case "agree":
		consent = true
	case "reject":
		consent = false
	default:
		return echo.ErrBadRequest
	}

	if consent {
		return authorizationIssueCaller(ctx)
	} else {
		return authorizationFailCaller(ctx, dto.AuthorizationFailReason_DENIED)
	}
}
