package handler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"mock-api/model"
	"net/http"
	"time"
)

type jwtCustomClaims struct {
	email string
	jwt.StandardClaims
}

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

	token, error := generateToken(identifier)
	if error != nil {
		return error
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func Login(c echo.Context) error {
	identifier := new(model.Identifier)
	if err := c.Bind(identifier); err != nil {
		return err
	}

	// FIXME Dummy
	if identifier.Email != "hoge@email.com" || identifier.Password != "aaa" {
		fmt.Println(identifier)
		return echo.ErrUnauthorized
	}

	token, error := generateToken(identifier)
	if error != nil {
		return error
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func generateToken(identifier *model.Identifier) (string, error) {
	claims := jwtCustomClaims{
		identifier.Email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, error := token.SignedString([]byte("secret"))
	if error != nil {
		return "", error
	}
	return accessToken, nil
}