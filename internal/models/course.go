package models


type Course struct {
	ID        int   `db:"id"`
	Name string `db:"name"`
}
