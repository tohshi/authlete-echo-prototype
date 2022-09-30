package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/authlete/authlete-go/api"
	"github.com/authlete/authlete-go/dto"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func authorizationHandler(ctx echo.Context) error {
	api := ctx.Get(AUTHLETE_API).(api.AuthleteApi)
	params := getParamsFromContext(ctx)

	authzRes, authleteErr := api.Authorization(&dto.AuthorizationRequest{Parameters: params})
	if authleteErr != nil {
		ctx.Logger().Error(authleteErr.Error())
		return echo.ErrInternalServerError
	}

	switch authzRes.Action {
	case dto.AuthorizationAction_INTERACTION:
		return authorizationInteractionHandler(ctx, authzRes)
	case dto.AuthorizationAction_BAD_REQUEST:
		return ctx.String(http.StatusBadRequest, authzRes.ResponseContent)
	case dto.AuthorizationAction_LOCATION:
		return ctx.Redirect(http.StatusFound, authzRes.ResponseContent)
	case dto.AuthorizationAction_INTERNAL_SERVER_ERROR:
		return ctx.String(http.StatusInternalServerError, authzRes.ResponseContent)
	default:
		return echo.ErrInternalServerError
	}
}

func authorizationInteractionHandler(ctx echo.Context, authzRes *dto.AuthorizationResponse) error {
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

	clientId := strconv.FormatUint(authzRes.Client.ClientId, 10)
	if authzRes.ClientIdAliasUsed {
		clientId = authzRes.Client.ClientIdAlias
	}

	authzSess.Values[TICKET] = authzRes.Ticket
	authzSess.Values[CLIENT_ID] = clientId
	scopesJson, _ := json.Marshal(authzRes.Scopes)
	authzSess.Values[SCOPES] = scopesJson
	authzSess.Options.MaxAge = 0

	if err := authzSess.Save(ctx.Request(), ctx.Response()); err != nil {
		ctx.Logger().Error(err)
		return err
	}

	if authnSess.Values[USER_ID] == nil {
		return ctx.Redirect(http.StatusFound, LOGIN_ENDPOINT)
	}

	if !isUserConsented(authnSess.Values[USER_ID].(string), clientId) {
		return ctx.Redirect(http.StatusFound, CONSENT_ENDPOINT)
	}

	return authorizationIssueCaller(ctx)
}

func authorizationIssueCaller(ctx echo.Context) error {
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

	api := ctx.Get(AUTHLETE_API).(api.AuthleteApi)
	authzIssueRes, authleteErr := api.AuthorizationIssue(&dto.AuthorizationIssueRequest{Ticket: authzSess.Values[TICKET].(string), Subject: authnSess.Values[USER_ID].(string)})
	if authleteErr != nil {
		return ctx.String(http.StatusBadRequest, authleteErr.Error())
	}

	return ctx.Redirect(http.StatusFound, authzIssueRes.ResponseContent)
}

func authorizationFailCaller(ctx echo.Context, reason dto.AuthorizationFailReason) error {
	sess, err := session.Get(AUTHORIZATION_SESSION, ctx)
	if err != nil {
		ctx.Logger().Error(err)
		return err
	}
	ticket := sess.Values[TICKET].(string)

	// Delete Authorization Session
	sess.Options.MaxAge = -1
	if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
		ctx.Logger().Error(err)
	}

	api := ctx.Get(AUTHLETE_API).(api.AuthleteApi)
	authzFailRes, authleteErr := api.AuthorizationFail(&dto.AuthorizationFailRequest{Ticket: ticket, Reason: reason})
	if authleteErr != nil {
		return ctx.String(http.StatusBadRequest, authleteErr.Error())
	}

	switch authzFailRes.Action {
	case dto.AuthorizationFailAction_INTERNAL_SERVER_ERROR:
		return ctx.String(http.StatusInternalServerError, authzFailRes.ResponseContent)
	case dto.AuthorizationFailAction_BAD_REQUEST:
		return ctx.String(http.StatusBadRequest, authzFailRes.ResponseContent)
	case dto.AuthorizationFailAction_LOCATION:
		return ctx.Redirect(http.StatusFound, authzFailRes.ResponseContent)
	default:
		return echo.ErrInternalServerError
	}
}
