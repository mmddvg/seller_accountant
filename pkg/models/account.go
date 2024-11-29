package models

type NewAccount struct {
	Name   string
	Charge uint
}

type Account struct {
	Id     uint   `db:"id"`
	Name   string `db:"name"`
	Charge uint   `db:"charge"`
}
