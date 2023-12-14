package entity

import "time"

type Endpoint struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Method      string `db:"method" json:"method"`
	Description string `db:"description" json:"description"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt time.Time `db:"deleted_at" json:"deletedAt"`
}
