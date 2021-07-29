package config

import (
	"jwt-vi-du-mau/model"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func ConnectDatabase() (db *pg.DB) {
	db = pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "postgres",
		Password: "123456",
		Database: "postgres",
	})

	err := createSchema(db)
	if err != nil {
		panic(err)
	}
	return db
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*model.User)(nil),
		(*model.Post)(nil),
		(*model.Comment)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
