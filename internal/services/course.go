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
		return []models.Course{}, fmt.Errorf("[in services.ListCourses] failed to get courses: %w", err)
	}

	defer rows.Close()

	var courses []models.Course

	for rows.Next(){
		var course models.Course
		err = rows.Scan(&course.ID, &course.Name)
		if err != nil {
			return []models.Course{}, fmt.Errorf("[in services.ListCourses] failed to scan user from row: %w", err)
		}
		courses = append(courses, course)

		if err = rows.Err(); err != nil {
			return []models.Course{}, fmt.Errorf("[in services.ListCourses] failed to scan courses: %w", err)
		}
	}

	return courses, nil
}

func (s CourseService) ListCourseByID(ctx context.Context, id int) (models.Course, error){
    row := s.database.QueryRowContext(
        ctx,
        `SELECT id, name FROM "course" WHERE id = $1`,
        id,
    )
	var course models.Course

	err := row.Scan(&course.ID, &course.Name)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Course{}, fmt.Errorf("[in services.GetCourseByID] no course found with ID: %w", err)
        }
        return models.Course{}, fmt.Errorf("[in services.GetCourseByID] failed to scan course: %w", err)
    }

	return course, nil
}

// UpdateCourse updates am CourseService objects from the database by ID.
func (s CourseService) UpdateCourse(ctx context.Context, ID int, course models.Course) (models.Course, error) {
	_, err := s.database.ExecContext(
		ctx,
		`
		UPDATE
			"course"
		SET
			"id" = $1,
			"name" = $2
		WHERE
			"id" = $3
		`,
		course.ID,
		course.Name,
		ID,
	)
	if err != nil {
		return models.Course{}, fmt.Errorf("[in services.UpdateUser] failed to update user: %w", err)
	}

	course.ID = int(ID)
	return course, nil
}

// CreateCourse updates am CourseService objects from the database by ID.
func (s CourseService) CreateCourse(ctx context.Context, course models.Course) (models.Course, error) {
	err:= s.database.QueryRowContext(
		ctx,
	   `INSERT INTO "course" (id, name)
		VALUES ($1, $2)
		RETURNING id, name;
		`,
		course.ID,
		course.Name,
	).Scan(&course.ID, &course.Name)

	if err != nil {
		return models.Course{}, fmt.Errorf("[in services.CreateCourse] failed to create course: %w", err)
	}

	return course, nil
}

func (s CourseService) DeleteCourseByID(ctx context.Context, id int) (models.Course, error){
	course, err := s.ListCourseByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			//check course first
			return models.Course{}, fmt.Errorf("[in services.DeleteCourseByID] no course found with ID: %d", id)
		}
		return models.Course{}, fmt.Errorf("[in services.DeleteCourseByID] failed to retrieve course: %w", err)
	}

	// Delete related entries from person_course table first
	_, err = s.database.ExecContext(
		ctx,
		`DELETE FROM "person_course" WHERE course_id = $1`,
		id,
	)
	if err != nil {
		return models.Course{}, fmt.Errorf("[in services.DeleteCourseByID] failed to delete from person_course: %w", err)
	}

	result, err := s.database.ExecContext(
		ctx,
		`DELETE FROM "course" WHERE id = $1`,
		id,
	)

	if err != nil {
		return models.Course{}, fmt.Errorf("[in services.DeleteCourseByID] failed to execute delete: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Course{}, fmt.Errorf("[in services.DeleteCourseByID] failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return models.Course{}, fmt.Errorf("[in services.DeleteCourseByID] no course found with ID: %d", id)
	}
	return course, nil

}