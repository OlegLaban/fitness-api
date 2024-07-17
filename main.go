package main

import (
	"github.com/labstack/echo/v4"
	"fitness-api/cmd/handlers"
	"fitness-api/cmd/storage"
)

func main() {
	e := echo.New()
	e.GET("/", handlers.Home)

	storage.InitDB()

	e.POST("/users", handlers.CreateUser)
	e.POST("/measurements", handlers.CreateMeasurement)
	e.PUT("/users/:id", handlers.HandleUpdateUser)
	e.PUT("/measurements/:id", handlers.HandleUpdateMeasurement)

	e.Logger.Fatal(e.Start(":8081"))
}