package handler

import (
	"github.com/labstack/echo"
	"mock-api/model"
	"net/http"
)

func Signup(c echo.Context) error {
	identifier := new(model.Identifier)
	if err := c.Bind(identifier); err != nil {
		return err
	}

	if identifier.Email == "" || identifier.Password == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid email or password",
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": "xxxxx",
	})
}
