package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/httplog/v2"
)

type ResponseMsg struct {
	Message string `json:"message"`
}

type ResponseErr struct {
	Error            string    `json:"error,omitempty"`
}

type OutputPerson struct {
	ID        int   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PersonType string `json:"type"`
	Age       int    `json:"age"`
}

type OutputCourse struct {
	ID        int   `json:"id"`
	Name string `json:"name"`
}

// MapOutput maps a models.User struct to an outputUser struct.
func MapOutput(person models.Person) OutputPerson {
	return OutputPerson{
		ID:        int(person.ID),
		FirstName: person.FirstName,
		LastName:  person.LastName,
		PersonType:      person.PersonType,
		Age:    person.Age,
	}
}

// MapOutput maps a models.User struct to an outputUser struct.
func MapOutputCourse(course models.Course) OutputCourse {
	return OutputCourse{
		ID:        int(course.ID),
		Name: course.Name,
	}
}

// MapMultipleOutput maps a slice of []models.User to a slice of []outputUser.
func MapMultipleOutput(person []models.Person) []OutputPerson {
	peopleOut := make([]OutputPerson, len(person))
	for i := 0; i < len(person); i++ {
		personOut := MapOutput(person[i])
		peopleOut[i] = personOut
	}

	return peopleOut
}

// MapMultipleOutput maps a slice of []models.User to a slice of []outputUser.
func MapMultipleCourseOutput(course []models.Course) []OutputCourse {
	coursesOut := make([]OutputCourse, len(course))
	for i := 0; i < len(course); i++ {
		courseOut := MapOutputCourse(course[i])
		coursesOut[i] = courseOut
	}

	return coursesOut
}


type ResponsePerson struct{
	Person OutputPerson `json:"person"`
}

type ResponsePeople struct {
	People [] OutputPerson `json:"people"`
}

type ResponseCourse struct{
	Course OutputCourse `json:"course"`
}

type ResponseCourses struct {
	Courses [] OutputCourse `json:"courses"`
}

// EncodeResponse encodes data as a JSON response.
func EncodeResponse(w http.ResponseWriter, logger *httplog.Logger, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Error while marshaling data", "err", err, "data", data)
		http.Error(w, `{"Error": "Internal server error"}`, http.StatusInternalServerError)
	}
}