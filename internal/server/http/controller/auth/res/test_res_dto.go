package httpctrl

type GetUserResponse struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
