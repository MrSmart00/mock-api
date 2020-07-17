package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"mock-api/handler"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	db := connectDB()
	defer db.Close()

	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)

	r := e.Group("/me")
	r.Use(middleware.JWT(handler.SigningKey()))
	r.POST("", handler.User)

	e.Logger.Fatal(e.Start(":3200"))
}

func connectDB() *gorm.DB {
	dbms := "mysql"
	user := "root"
	password := "password"
	protocol := "tcp(dummy-mysql:3306)"
	dbname := "mysql"

	connect := user + ":" + password + "@" + protocol + "/" + dbname
	db, error := gorm.Open(dbms, connect)
	if error != nil {
		panic(error.Error())
	}
	return db
}
