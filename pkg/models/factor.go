package models

type NewFactor struct {
	Products  []uint
	AccountId uint
}

type Factor struct {
	Id        uint `db:"id"`
	Products  []uint
	AccountId uint `db:"account_id"`
}
