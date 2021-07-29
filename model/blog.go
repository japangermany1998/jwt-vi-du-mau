package model

import "time"

type Post struct{
	tableName struct{} `pg:"blog.post"`

	Id int `pg:"type:serial" `

	Content string `pg:",notnull"`

	Title string	`pg:",notnull"`

	CreatedAt time.Time `pg:"type:timestamp without time zone,default:now()"`

	UpdatedAt time.Time `pg:"type:timestamp without time zone"`

	UserId int `pg:"type:integer"`
	
	User *User `pg:"rel:has-one"`
}