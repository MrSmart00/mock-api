package handler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"mock-api/model"
	"net/http"
	"time"
)

type JwtCustomClaims struct {
	email string
	jwt.StandardClaims
}

func SigningKey() []byte {
	return []byte("secret")
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
	uuid, error := uuid.NewUUID()
	if error != nil {
		return "", error
	}
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = uuid.String()
	claims["email"] = identifier.Email
	claims["expired"] = time.Now().Add(time.Hour * 2).Unix()
	accessToken, error := token.SignedString(SigningKey())
	if error != nil {
		return "", error
	}
	return accessToken, nil
}

func User(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uuid := claims["uuid"].(string)
	email := claims["email"].(string)
	expired := time.Unix(int64(claims["expired"].(float64)), 0)
	return c.JSON(http.StatusOK, echo.Map {
		"message": "Hello "+email,
		"id": uuid,
		"expired_date": expired,
	})
}
