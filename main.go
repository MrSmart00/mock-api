package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"mock-api/handler"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)

	r := e.Group("/user")
	r.Use(middleware.JWT([]byte(handler.SigningKey)))
	r.GET("", handler.Restricted)

	e.Logger.Fatal(e.Start(":3200"))
}
