package router

import (
	"jwt-vi-du-mau/controller"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	app.Post("/user/login", controller.Login)

	app.Get("/user", controller.User)

	app.Get("/logout", controller.Logout)

	app.Post("/register", controller.Register)

}
