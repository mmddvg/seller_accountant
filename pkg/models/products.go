package models

type NewProduct struct {
	Name  string
	Price uint
}

type Product struct {
	Id    uint   `db:"id"`
	Name  string `db:"name"`
	Price uint   `db:"price"`
}
