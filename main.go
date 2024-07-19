package main

import (
	"errors"
	"fitness-api/cmd/handlers"
	"fitness-api/cmd/storage"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", handlers.Home)
	
	storage.InitDB()
	e.POST("/token", handlers.Auth)
	e.POST("/users", handlers.CreateUser)
	e.POST("/measurements", handlers.CreateMeasurement, auth)
	e.PUT("/users/:id", handlers.HandleUpdateUser, auth)
	e.PUT("/measurements/:id", handlers.HandleUpdateMeasurement, auth)
	e.GET("/users", handlers.HandleGetUsers, auth)

	e.Logger.Fatal(e.Start(":8081"))
}

func auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Auth")
		err := handlers.ValidateToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, errors.New("Token is invalid").Error())
		}
		
		return next(c)
	}
}