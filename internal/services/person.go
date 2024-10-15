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
		return []models.Person{}, fmt.Errorf("[in services.ListPeople] failed to get people: %w", err)
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
			return []models.Person{}, fmt.Errorf("[in services.ListPeople] failed to scan people: %w", err)
		}
	}

	return people, nil
}

func (s PersonService) ListPersonByID(ctx context.Context, id int) (models.Person, error){
    row := s.database.QueryRowContext(
        ctx,
        `SELECT id, first_name, last_name, type, age FROM "person" WHERE id = $1`,
        id,
    )
	var person models.Person

	err := row.Scan(&person.ID, &person.FirstName, &person.LastName, &person.PersonType, &person.Age)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Person{}, fmt.Errorf("[in services.GetPersonByID] no person found with ID: %w", err)
        }
        return models.Person{}, fmt.Errorf("[in services.GetPersonByID] failed to scan person: %w", err)
    }

	return person, nil
}

// CreatePerson creates a PersonService objects from the database by ID.
func (s PersonService) CreatePerson(ctx context.Context, person models.Person) (models.Person, error) {
	err:= s.database.QueryRowContext(
		ctx,
	   `INSERT INTO "person" (id, first_name, last_name, type, age)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, first_name, last_name, type, age;
		`,
		person.ID,
		person.FirstName,
		person.LastName,
		person.PersonType,
		person.Age,
	).Scan(&person.ID, &person.FirstName, &person.LastName, &person.PersonType, &person.Age)

	if err != nil {
		return models.Person{}, fmt.Errorf("[in services.CreatePerson] failed to create person: %w", err)
	}

	return person, nil
}

func (s PersonService) DeletePersonByID(ctx context.Context, id int) (models.Person, error){
	person, err := s.ListPersonByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			//check person first
			return models.Person{}, fmt.Errorf("[in services.DeletePersonByID] no person found with ID: %d", id)
		}
		return models.Person{}, fmt.Errorf("[in services.DeletePersonByID] failed to retrieve person: %w", err)
	}

	_, err = s.database.ExecContext(
		ctx,
		`DELETE FROM "person_course" WHERE person_id = $1`,
		id,
	)
	if err != nil {
		return models.Person{}, fmt.Errorf("[in services.DeletePersonByID] failed to delete from person_course: %w", err)
	}

	result, err := s.database.ExecContext(
		ctx,
		`DELETE FROM "person" WHERE id = $1`,
		id,
	)

	if err != nil {
		return models.Person{}, fmt.Errorf("[in services.DeletePersonByID] failed to execute delete: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Person{}, fmt.Errorf("[in services.DeletePersonByID] failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return models.Person{}, fmt.Errorf("[in services.DeletePersonByID] no person found with ID: %d", id)
	}
	return person, nil

}

// UpdatePerson updates am UserService objects from the database by ID.
func (s PersonService) UpdatePerson(ctx context.Context, ID int, person models.Person) (models.Person, error) {
	_, err := s.database.ExecContext(
		ctx,
		`
		UPDATE
			"person"
		SET
			"id" = $1,
			"first_name" = $2,
			"last_name" = $3,
			"type" = $4,
			"age"= $5
		WHERE
			"id" = $6
		`,
		person.ID,
		person.FirstName,
		person.LastName,
		person.PersonType,
		person.Age,
		ID,
	)
	if err != nil {
		return models.Person{}, fmt.Errorf("[in services.UpdateUser] failed to update user: %w", err)
	}

	person.ID = int(ID)
	return person, nil
}