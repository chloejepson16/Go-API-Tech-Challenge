package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
)

type InputPerson struct {
	ID        int   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PersonType string `json:"type"`
	Age       int    `json:"age"`
}

type InputCourse struct{
	ID        int   `json:"id"`
	Name string `json:"name"`
}

// MapTo maps a inputUser to a models.User object.
func (person InputPerson) MapTo() (models.Person, error) {
	return models.Person{
		ID:        int(person.ID),
		FirstName: person.FirstName,
		LastName:  person.LastName,
		PersonType:      person.PersonType,
		Age:    person.Age,
	}, nil
}

// MapTo maps a inputUser to a models.User object.
func (course InputCourse) MapTo() (models.Course, error) {
	return models.Course{
		ID:        int(course.ID),
		Name: course.Name,
	}, nil
}

// Valid validates all fields of an inputUser struct.
func (user InputPerson) Valid() []problem {
	var problems []problem

	// validate FirstName is not blank
	if user.FirstName == "" {
		problems = append(problems, problem{
			Name:        "first_name",
			Description: "must not be blank",
		})
	}

	// validate LastName is not blank
	if user.LastName == "" {
		problems = append(problems, problem{
			Name:        "last_name",
			Description: "must not be blank",
		})
	}

	if user.PersonType == "" {
		problems = append(problems, problem{
			Name:        "ID",
			Description: "must not be blank",
		})
	}

	// validate UserID greater than 0
	if user.ID < 1 {
		problems = append(problems, problem{
			Name:        "user_id",
			Description: "must be must be greater than zero",
		})
	}

	return problems
}

// Valid validates all fields of an inputUser struct.
func (course InputCourse) Valid() []problem {
	var problems []problem

	// validate LastName is not blank
	if course.ID < 0 {
		problems = append(problems, problem{
			Name:        "id",
			Description: "must not be less than 0",
		})
	}

	if course.Name == "" {
		problems = append(problems, problem{
			Name:        "name",
			Description: "must not be blank",
		})
	}

	return problems
}

// problem represents an issue found during validation.
type problem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Validator is an interface that defines a method for validating an object.
// It returns a slice of problems found during validation.
type Validator interface {
	Valid() (problems []problem)
}

// Mapper is a generic interface that defines a method for mapping an object to another type.
// The MapTo method returns the mapped object and an error if the mapping fails.
type Mapper[T any] interface {
	MapTo() (T, error)
}

// ValidatorMapper is a generic interface that combines Validator and Mapper interfaces.
// It requires implementing both validation and mapping methods.

type ValidatorMapper[T any] interface {
	Validator
	Mapper[T]
}

// DecodeValidateBody decodes a JSON string into a ValidatorMapper, validates it, and maps it to
// the output type. If decoding, validation, or mapping fails, it returns the appropriate errors
// and problems.
func DecodeValidateBody[I ValidatorMapper[O], O any](r *http.Request) (O, []problem, error) {
	var inputModel I

	// decode to JSON
	if err := json.NewDecoder(r.Body).Decode(&inputModel); err != nil {
		return *new(O), nil, fmt.Errorf("[in DecodeValidateBody] decode json: %w", err)
	}

	// validate
	if problems := inputModel.Valid(); len(problems) > 0 {
		return *new(O), problems, fmt.Errorf(
			"[in DecodeValidateBody] invalid %T: %d problems", inputModel, len(problems),
		)
	}

	// map to return type
	data, err := inputModel.MapTo()
	if err != nil {
		return *new(O), nil, fmt.Errorf(
			"[in DecodeValidateBody] error mapping input %T to %T: %w",
			*new(I),
			*new(O),
			err,
		)
	}

	return data, nil, nil
}
