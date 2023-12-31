package entity

import (
	"time"
)

type AppUser struct {
	Id                  int64                    `db:"id" json:"id"`
	Username            string                   `db:"username" json:"username"`
	Email               string                   `db:"email" json:"email"`
	Password            string                   `db:"password" json:"password"`
	FirstName           string                   `db:"first_name" json:"firstName"`
	LastName            string                   `db:"last_name" json:"lastName"`
	PhoneNumber         string                   `db:"phone_number" json:"phoneNumber"`
	EmailConfirmed      bool                     `db:"email_confirmed" json:"emailConfirmed"`
	FpEmailConfirmation EmailConfirmationDetails `db:"fp_email_confirmation" json:"fpEmailConfirmationDetails"`
	EmailConfirmation   EmailConfirmationDetails `db:"email_confirmation" json:"emailConfirmationDetails"`
	RoleIds             []int                    `db:"role_ids" json:"roles"`
	BlockDetails        []BlockDetails           `db:"block_details" json:"blockDetails"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt time.Time `db:"deleted_at" json:"deletedAt"`
}

type EmailConfirmationDetails struct {
	Code      string    `db:"code" json:"code"`
	ExpiresAt time.Time `db:"expires_at" json:"expiresAt"`
}

type BlockDetails struct {
	Reason    string    `db:"reason" json:"reason" validate:"required"`
	UntilAt   time.Time `db:"until_at" json:"untilAt" validate:"required,gtNow"`
	BlockedBy string    `db:"blocked_by" json:"blockedBy"`
}
