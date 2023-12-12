package app_user

import "time"

type AppUser struct {
	Id                  int64                    `db:"id" json:"id"`
	Username            string                   `db:"username" json:"username"`
	Email               string                   `db:"email" json:"email"`
	Password            string                   `db:"password" json:"password"`
	EmailConfirmed      bool                     `db:"email_confirmed" json:"emailConfirmed"`
	FpEmailConfirmation EmailConfirmationDetails `db:"fp_email_confirmation" json:"fpEmailConfirmationDetails"`
	EmailConfirmation   EmailConfirmationDetails `db:"email_confirmation" json:"emailConfirmationDetails"`
	Roles               []int                    `db:"roles" json:"roles"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt time.Time `db:"deleted_at" json:"deletedAt"`
}

type EmailConfirmationDetails struct {
	Code      string    `db:"code" json:"code"`
	ExpiresAt time.Time `db:"expires_at" json:"expiresAt"`
}
