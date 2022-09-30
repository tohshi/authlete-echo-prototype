package main

import (
	"github.com/authlete/authlete-go/api"
	"github.com/labstack/echo/v4"
)

type AuthleteApi struct {
	api.AuthleteApi
}

func (api *AuthleteApi) authleteMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Set(AUTHLETE_API, api)
		return next(ctx)
	}
}
