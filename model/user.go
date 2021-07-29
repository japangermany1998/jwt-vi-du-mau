package model

type User struct{
	tableName struct{} `pg:"auth.users"`

	Id int `pg:"type:serial"`

	FirstName string

	LastName string

	Email string	`pg:",unique"`

	Password string	`json:"-"`
	
	Posts []*Post `pg:"rel:has-many"` 
}