package main

import (
	"jwt-vi-du-mau/config"
	"jwt-vi-du-mau/controller"
	"jwt-vi-du-mau/middleware"
	"jwt-vi-du-mau/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db := config.ConnectDatabase() //Kết nối database
	defer db.Close()

	app := fiber.New()

	controller.DB = db

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	api := app.Group("/api").Use(middleware.IsAuthenticated)
	router.AllRoutes(api) //tất cả các route trong này sẽ phải check authenticate trước

	router.AuthRoutes(app)

	app.Listen(":8080")

}
