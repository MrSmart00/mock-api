package handler

import (
	"mock-api/db"
	"mock-api/model"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type JwtCustomClaims struct {
	email string
	jwt.StandardClaims
}

func SigningKey() []byte {
	return []byte("secret")
}

type Auth struct {
	DB *db.ImplDB
}

func New(db *db.ImplDB) *Auth {
	auth := new(Auth)
	auth.DB = db
	return auth
}

func (auth Auth) Signup(context echo.Context) error {
	identifier := new(model.Identifier)
	if err := context.Bind(identifier); err != nil {
		return err
	}

	if identifier.Email == "" || identifier.Password == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid email or password",
		}
	}

	token := jwt.New(jwt.SigningMethodHS256)
	uuid, error := uuid.NewUUID()
	if error != nil {
		return error
	}
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = uuid.String()
	claims["email"] = identifier.Email
	claims["expired"] = time.Now().Add(time.Minute * 5).Unix()
	accessToken, error := token.SignedString(SigningKey())
	if error != nil {
		return error
	}

	user := db.User{}
	user.AccessToken = accessToken
	user.UserID = uuid.String()
	user.Email = identifier.Email
	user.Password = identifier.Password
	user.CreatedAt = time.Now()
	user.LoggedInAt = time.Now()
	auth.DB.Create(user)

	return context.JSON(http.StatusOK, echo.Map{
		"token": accessToken,
	})
}

func (auth Auth) Login(context echo.Context) error {
	identifier := new(model.Identifier)
	if err := context.Bind(identifier); err != nil {
		return err
	}

	user, error := auth.DB.Find(*identifier)
	if error != nil {
		return error
	}
	user = auth.DB.UpdateLoggedInAt(user)

	return context.JSON(http.StatusOK, echo.Map{
		"token": user.AccessToken,
	})
}

func (auth Auth) User(context echo.Context) error {
	user := context.Get("user").(*jwt.Token)
	accessToken, error := user.SignedString(SigningKey())
	if error != nil {
		return error
	}
	foundUser, error := auth.DB.FindByToken(accessToken)
	if error != nil {
		return error
	}
	return context.JSON(http.StatusOK, echo.Map{
		"email":      foundUser.Email,
		"uuid":       foundUser.UserID,
		"created_at": foundUser.CreatedAt,
		"logged_in_at": 	foundUser.LoggedInAt,
	})
}

func (auth Auth) DeleteAccount(context echo.Context) error {
	user := context.Get("user").(*jwt.Token)
	accessToken, error := user.SignedString(SigningKey())
	if error != nil {
		return error
	}
	if err := auth.DB.DeleteByToken(accessToken); err != nil {
		return err
	}
	return context.NoContent(http.StatusOK)
}
