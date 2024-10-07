package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
)

type PersonService struct{
	database *sql.DB
}

//NewPersonService will return a NewPerson struct,

func NewPersonService(db *sql.DB) *PersonService{
	return &PersonService{database: db,}
}

func (s PersonService) ListPeople(ctx context.Context) ([]models.Person, error){
	rows, err:= s.database.QueryContext(
		ctx,
		`SELECT * FROM "person"`,
	)

	if err != nil{
		return []models.Person{}, fmt.Errorf("[in services.ListPeople] failed to get users: %w", err)
	}

	defer rows.Close()

	var people []models.Person

	for rows.Next(){
		var person models.Person
		err = rows.Scan(&person.ID, &person.FirstName, &person.LastName, &person.PersonType, &person.Age)
		if err != nil {
			return []models.Person{}, fmt.Errorf("[in services.ListPeople] failed to scan user from row: %w", err)
		}
		people = append(people, person)

		if err = rows.Err(); err != nil {
			return []models.Person{}, fmt.Errorf("[in services.ListPeople] failed to scan users: %w", err)
		}
	}

	return people, nil
}
