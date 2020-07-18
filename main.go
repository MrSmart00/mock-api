package main

import (
	"mock-api/db"
	"mock-api/handler"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	db := &db.ImplDB{}
	db.Start()
	defer db.Close()


	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world!")
	})

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	auth := handler.New(db)

	e.POST("/signup", auth.Signup)
	e.POST("/login", auth.Login)

	r := e.Group("/me")
	r.Use(middleware.JWT(handler.SigningKey()))
	r.POST("", auth.User)
	r.DELETE("/delete", auth.DeleteAccount)

	e.Logger.Fatal(e.Start(":3200"))
}
