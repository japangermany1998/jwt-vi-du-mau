package controller

import (
	"jwt-vi-du-mau/model"
	"jwt-vi-du-mau/util"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func User(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("jwt")	//cookie chứa value là token đăng nhập

	issuer, _ := util.ParseJWT(cookie)	//Parse token lấy thông tin đăng nhập

	var user model.User

	DB.Model(&user).Where("id = ?", issuer).Relation("Posts").Select()

	return ctx.JSON(user)
}

func Logout(ctx *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)

	return ctx.JSON(fiber.Map{
		"message": "Logout success",
	})
}

func Login(ctx *fiber.Ctx) error {
	var data map[string]string
	if err := ctx.BodyParser(&data); err != nil {
		return ctx.JSON(err)
	}

	var user model.User
	err := DB.Model(&user).Where("email = ?", data["email"]).First()

	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(err)
	}

	if user.Id == 0 {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(fiber.Map{
			"message": "not found",
		})
	}

	//So sánh password nhập vào với password băm trong database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		ctx.Status(400)

		return ctx.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	token, err := util.GenerateJWT(strconv.Itoa(user.Id)) //Tạo token đăng nhập
	if err != nil {
		log.Print(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{ //Cookie có value là token
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24), //thời hạn 1 ngày
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)

	return ctx.JSON(token)
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
