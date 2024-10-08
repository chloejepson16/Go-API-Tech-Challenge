package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/models"
)

type CourseService struct{
	database *sql.DB
}

func NewCourseService(db *sql.DB) *CourseService{
	return &CourseService{database: db,}
}

func (s CourseService) ListCourses(ctx context.Context) ([]models.Course, error){
	rows, err:= s.database.QueryContext(
		ctx,
		`SELECT * FROM "course"`,
	)

	if err != nil{
		return []models.Course{}, fmt.Errorf("[in services.ListPeople] failed to get people: %w", err)
	}

	defer rows.Close()

	var courses []models.Course

	for rows.Next(){
		var course models.Course
		err = rows.Scan(&course.ID, &course.Name)
		if err != nil {
			return []models.Course{}, fmt.Errorf("[in services.ListPeople] failed to scan user from row: %w", err)
		}
		courses = append(courses, course)

		if err = rows.Err(); err != nil {
			return []models.Course{}, fmt.Errorf("[in services.ListPeople] failed to scan people: %w", err)
		}
	}

	return courses, nil
}