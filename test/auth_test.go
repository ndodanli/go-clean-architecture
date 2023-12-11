package test

import (
	"testing"
)

func TestLogin(t *testing.T) {
	defer setupTest()()

	//tableTest := []struct {
	//	name    string
	//	payload *req.LoginRequest
	//	want    string
	//}{
	//	{"fail authenticate", &req.LoginRequest{Username: "test", Password: "test1234"}, "Username or password is incorrect"},
	//	{"success authenticate", &req.LoginRequest{Username: "test", Password: "test123"}, ""},
	//}
	//
	//for _, param := range tableTest {
	//	t.Run(param.name, func(t *testing.T) {
	//		res := appServices.AuthService.Login(ctx, *param.payload, ts)
	//		got := res.GetErrorMessage()
	//		assert.DeepEqual(t, got, param.want)
	//	})
	//}

}
