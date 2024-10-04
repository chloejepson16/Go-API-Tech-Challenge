package services

import (
	"database/sql"
)

type PersonService struct{
	database *sql.DB
}

//NewPersonService will return a NewPerson struct,

func NewPersonService(db *sql.DB) *PersonService{
	return &PersonService{database: db,}
}
