package main

import (
	"net/http"

	"github.com/authlete/authlete-go/dto"
	"github.com/labstack/echo/v4"
)

func tokenHandler(ctx echo.Context) error {
	username, password, _ := ctx.Request().BasicAuth()
	params := getParamsFromContext(ctx)

	tokenRes, authleteErr := authleteApi.Token(&dto.TokenRequest{Parameters: params, ClientId: username, ClientSecret: password})
	if authleteErr != nil {
		ctx.Logger().Error(authleteErr.Error())
		return echo.ErrInternalServerError
	}

	switch tokenRes.Action {
	case dto.TokenAction_OK:
		return ctx.String(http.StatusOK, tokenRes.ResponseContent)
	case dto.TokenAction_INVALID_CLIENT:
		return ctx.String(http.StatusUnauthorized, tokenRes.ResponseContent)
	case dto.TokenAction_INTERNAL_SERVER_ERROR:
		return ctx.String(http.StatusInternalServerError, tokenRes.ResponseContent)
	case dto.TokenAction_BAD_REQUEST:
		return ctx.String(http.StatusBadRequest, tokenRes.ResponseContent)
	default:
		return echo.ErrInternalServerError
	}
}
