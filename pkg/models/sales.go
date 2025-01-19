package models

type Sale struct {
	Id         uint `db:"id"`
	CustomerId uint `db:"customer_id"`
	Price      int  `db:"price"`
}
