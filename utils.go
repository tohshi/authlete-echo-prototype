package main

import "github.com/labstack/echo/v4"

func getParamsFromContext(ctx echo.Context) string {
	var params string

	switch ctx.Request().Method {
	case "GET":
		params = ctx.Request().URL.RawQuery
	case "POST":
		p, err := ctx.FormParams()
		if err == nil {
			params = p.Encode()
		}
	}

	return params
}
