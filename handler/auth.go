package handler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"mock-api/model"
	"net/http"
	"time"
)

const SigningKey = "secret"

type JwtCustomClaims struct {
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

	return c.JSON(http.StatusOK, echo.Map{
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

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func generateToken(identifier *model.Identifier) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = identifier.Email
	claims["expired"] = time.Now().Add(time.Hour * 2).Unix()
	accessToken, error := token.SignedString([]byte(SigningKey))
	if error != nil {
		return "", error
	}
	return accessToken, nil
}

func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	expired := claims["expired"].(float64)
	fmt.Println(time.Unix(int64(expired), 0))
	return c.String(http.StatusOK, "Welcome "+email+"!")
}
