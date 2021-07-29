package controller

import (
	"jwt-vi-du-mau/model"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetComments(ctx *fiber.Ctx) error {
	var comments []model.Comment
	err := DB.Model(&comments).Relation("User").Relation("Post").Select()
	if err != nil {
		panic(err)
	}
	return ctx.JSON(comments)
}

func GetCommentById(ctx *fiber.Ctx) error {
	id := ctx.Params("userId")

	var comment model.Comment

	err := DB.Model(&comment).Relation("Post").Relation("User").Where("id = ?", id).Select()
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}

	return ctx.JSON(comment)

}

func UpdateComment(ctx *fiber.Ctx) error {
	var data map[string]interface{}
	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	id := ctx.Params("id")

	_, err := DB.Model(&data).TableExpr("blog.comment").Where("id = ?", id).Update()
	if err != nil {
		log.Println(err)
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}

	return ctx.JSON("Update thành công")
}

func DeleteComment(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	_, err := DB.Model((*model.Comment)(nil)).Where("id = ?", id).Delete()
	if err != nil {
		log.Println(err)
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}
	return ctx.JSON("Xóa thành công")
}

func CreateComment(ctx *fiber.Ctx) error {
	var data map[string]interface{}
	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	data["user_id"] = 1
	data["post_id"] = 1

	_, err := DB.Model(&data).TableExpr("blog.comment").Insert()
	if err != nil {
		log.Println(err)
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}

	return ctx.JSON("Comment thành công")
}
