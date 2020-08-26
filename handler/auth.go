package handler

import (
	"github.com/labstack/echo/middleware"
	"mock-api/db"
	"mock-api/model"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type jwtCustomClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

func signingKey() []byte {
	return []byte("secret")
}

func JWTConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: signingKey(),
	}
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

	uuid, error := uuid.NewUUID()
	if error != nil {
		return error
	}
	uuidString := uuid.String()
	claims := &jwtCustomClaims{
		uuidString,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, error := token.SignedString(signingKey())
	if error != nil {
		return error
	}

	user := db.User{}
	user.AccessToken = accessToken
	user.UserID = uuidString
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
	token := context.Get("user").(*jwt.Token)
	claims := token.Claims.(*jwtCustomClaims)
	foundUser, error := auth.DB.FindByUserId(claims.ID)
	if error != nil {
		return error
	}
	return context.JSON(http.StatusOK, echo.Map{
		"email":      	foundUser.Email,
		"id":       		foundUser.UserID,
		"created_at": 	foundUser.CreatedAt.Format(time.RFC3339),
		"logged_in_at": foundUser.LoggedInAt.Format(time.RFC3339),
	})
}

func (auth Auth) DeleteAccount(context echo.Context) error {
	user := context.Get("user").(*jwt.Token)
	accessToken, error := user.SignedString(signingKey())
	if error != nil {
		return error
	}
	if err := auth.DB.DeleteByToken(accessToken); err != nil {
		return err
	}
	return context.NoContent(http.StatusOK)
}
