package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type appError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := "internal server error"

	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		switch msg := he.Message.(type) {
		case string:
			message = msg
		case error:
			message = msg.Error()
		default:
			message = http.StatusText(code)
		}
	}

	c.Logger().Errorf("HTTP %d: %v", code, err)

	_ = c.JSON(code, appError{
		Code:    code,
		Message: message,
	})
}
