package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/samber/lo"
)

type Factor struct {
	ID         int            `db:"id"`
	PurchaseID int            `db:"purchase_id"`
	StoreName  string         `db:"store_name"`
	Price      int            `db:"price"`
	FileName   sql.NullString `db:"file_name"`
}

func (f *Factor) String() string {
	return fmt.Sprintf("name : %s , price : %d ", f.StoreName, f.Price)
}

type Purchase struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Factors   []Factor
}

func (p *Purchase) String() string {
	return fmt.Sprintf("id : %d , factors : %s , created at : %s", p.ID, lo.Reduce(p.Factors, func(agg string, f Factor, i int) string { return agg + " " + f.String() }, "[ ")+" ]", p.CreatedAt.String())
}
