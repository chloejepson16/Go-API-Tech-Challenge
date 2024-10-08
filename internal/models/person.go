package models


type Person struct {
	ID        int   `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	PersonType string `db:"type"`
	Age       int    `db:"age"`
}
