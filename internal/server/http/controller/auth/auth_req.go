package authctrl

type GetUserRequest struct {
	ID int `query:"ID" validate:"required,numeric,min=10,max=20"`
}
