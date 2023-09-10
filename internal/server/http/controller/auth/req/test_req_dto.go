package httpctrl

type GetUserRequest struct {
	ID int `query:"id" validate:"required,numeric,min=10,max=20"`
}
