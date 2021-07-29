package model

import "time"

type Comment struct{
	tableName struct{} `pg:"blog.comment"`

	Id int `pg:"type:serial"`

	Content string	`pg:",notnull"`

	CreatedAt time.Time	`pg:"type:timestamp without time zone,default:now()"`

	UserId int	`pg:"type:integer,notnull"`

	User *User `pg:"rel:has-one"`

	PostId int	`pg:"type:integer,notnull"`
	
	Post *Post	`pg:"rel:has-one"`
}