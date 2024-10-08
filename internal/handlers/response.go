package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/httplog/v2"
)

type responseMsg struct {
	Message string `json:"message"`
}

type responseErr struct {
	Error            string    `json:"error,omitempty"`
}

type outputPerson struct {
	ID        int   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PersonType string `json:"type"`
	Age       int    `json:"age"`
}

type outputCourse struct {
	ID        int   `json:"id"`
	Name string `json:"name"`
}

// mapOutput maps a models.User struct to an outputUser struct.
func mapOutput(person models.Person) outputPerson {
	return outputPerson{
		ID:        int(person.ID),
		FirstName: person.FirstName,
		LastName:  person.LastName,
		PersonType:      person.PersonType,
		Age:    person.Age,
	}
}

// mapOutput maps a models.User struct to an outputUser struct.
func mapOutputCourse(course models.Course) outputCourse {
	return outputCourse{
		ID:        int(course.ID),
		Name: course.Name,
	}
}

// mapMultipleOutput maps a slice of []models.User to a slice of []outputUser.
func mapMultipleOutput(person []models.Person) []outputPerson {
	peopleOut := make([]outputPerson, len(person))
	for i := 0; i < len(person); i++ {
		personOut := mapOutput(person[i])
		peopleOut[i] = personOut
	}

	return peopleOut
}

// mapMultipleOutput maps a slice of []models.User to a slice of []outputUser.
func mapMultipleCourseOutput(course []models.Course) []outputCourse {
	coursesOut := make([]outputCourse, len(course))
	for i := 0; i < len(course); i++ {
		courseOut := mapOutputCourse(course[i])
		coursesOut[i] = courseOut
	}

	return coursesOut
}


type responsePerson struct{
	Person outputPerson `json:"person"`
}

type responsePeople struct {
	People [] outputPerson `json:"people"`
}

type responseCourse struct{
	Course outputCourse `json:"course"`
}

type responseCourses struct {
	Courses [] outputCourse `json:"courses"`
}

// encodeResponse encodes data as a JSON response.
func encodeResponse(w http.ResponseWriter, logger *httplog.Logger, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Error while marshaling data", "err", err, "data", data)
		http.Error(w, `{"Error": "Internal server error"}`, http.StatusInternalServerError)
	}
}