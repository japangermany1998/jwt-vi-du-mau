package controller

import (
	"jwt-vi-du-mau/model"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetPosts(ctx *fiber.Ctx) error {
	var posts []model.Post
	err := DB.Model(&posts).Relation("User").Select()
	if err != nil {
		panic(err)
	}
	return ctx.JSON(posts)
}

func GetPostById(ctx *fiber.Ctx) error {
	id := ctx.Params("userId")

	var post model.Post

	err := DB.Model(&post).Relation("User").Where("id = ?", id).Select()
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}

	return ctx.JSON(post)

}

func CreatePost(ctx *fiber.Ctx) error {
	var data map[string]interface{}
	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	data["user_id"] = 1

	_, err := DB.Model(&data).TableExpr("blog.post").Insert()
	if err != nil {
		log.Println(err)
		ctx.Status(500)
		return err
	}

	return ctx.JSON("Tạo bài viết thành công")
}

func UpdatePost(ctx *fiber.Ctx) error {
	var data map[string]interface{}
	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	id := ctx.Params("id")

	data["updated_at"] = time.Now()
	_, err := DB.Model(&data).TableExpr("blog.post").Where("id = ?", id).Update()
	if err != nil {
		log.Println(err)
		ctx.Status(500)
		return err
	}

	return ctx.JSON("Cập nhật thành công")
}

func DeletePost(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	_, err := DB.Model((*model.Post)(nil)).Where("id = ?", id).Delete()
	if err != nil {
		log.Println(err)
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}
	return ctx.JSON("Xóa thành công")
}
