package models

type Customer struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	Charge int    `db:"charge"`
}
