package router

import (
	"jwt-vi-du-mau/controller"

	"github.com/gofiber/fiber/v2"
)

func AllRoutes(app fiber.Router) {
	app.Get("/users", controller.GetUsers)

	app.Get("/user/{userId}", controller.GetUserById)

	app.Put("/user/{id}", controller.UpdateUser)

	app.Get("/posts", controller.GetPosts)

	app.Get("/posts/{postId}", controller.GetPostById)

	app.Post("/post/create", controller.CreatePost)

	app.Put("/post/{id}", controller.UpdatePost)

	app.Delete("/post/{id}", controller.DeletePost)

	app.Get("/comment", controller.GetComments)

	app.Get("/comment/{commentId}", controller.GetCommentById)

	app.Post("/comment/create", controller.CreateComment)

	app.Put("/comment/{id}", controller.UpdateComment)

	app.Delete("/comment/{id}", controller.DeleteComment)

}
