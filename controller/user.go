package controller

import (
	"jwt-vi-du-mau/model"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var DB *pg.DB

func GetUsers(ctx *fiber.Ctx) error {
	var user []model.User

	err := DB.Model(&user).Relation("Posts").Select()

	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}

	return ctx.JSON(user)
}

func GetUserById(ctx *fiber.Ctx) error {
	id := ctx.Params("userId")

	var user model.User

	err := DB.Model(&user).Relation("Posts").Where("id = ?", id).Select()
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}

	return ctx.JSON(user)
}

func Register(ctx *fiber.Ctx) error {
	var data map[string]string

	if err := ctx.BodyParser(&data); err != nil {
		log.Print(err)
		return err
	}

	if data["password"] != data["passwordconfirm"] {
		ctx.Status(400)
		return ctx.JSON(map[string]string{
			"message": "password doesn't match",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := model.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  string(password),
	}

	_, err := DB.Model(&user).Insert()
	if err != nil {
		panic(err)
	}
	return ctx.JSON(user)
}

func UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var data map[string]interface{}
	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	_, err := DB.Model(&data).TableExpr("auth.users").Where("id = ?", id).Update()
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}

	ctx.Status(fiber.StatusOK)
	return ctx.JSON("Cập nhật thành công")
}
