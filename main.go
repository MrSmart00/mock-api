package main

import (
	"github.com/labstack/echo"
	"mock-api/handler"
)

func main() {
	e := echo.New()

	e.POST("/signup", handler.Signup)

	e.Logger.Fatal(e.Start(":3200"))
}
