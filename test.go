// ForgotPassword godoc
// @Security BearerAuth
// @Summary      ForgotPassword
// @Description  ForgotPassword
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        ForgotPassword path string true "Refresh Token"
// @Success      200  {object}   baseres.SwaggerSuccessRes[queries.ForgotPasswordQueryResponse] "OK. On success."
// @Failure      400  {object}   baseres.SwaggerValidationErrRes "Bad Request. On any validation error."
// @Failure      401  {object}   baseres.SwaggerUnauthorizedErrRes "Unauthorized."
// @Failure      500  {object}   baseres.SwaggerInternalErrRes "Internal Server Error."
// @Router       /v1/auth/ForgotPassword [get]
func (ac *AuthController) ForgotPassword(c echo.Context) error {
var query queries.ForgotPasswordQuery
if err := utils.BindAndValidate(c, &query); err != nil {
return err
}
res := mediatr.Send[*queries.ForgotPasswordQuery, *baseres.Result[queries.ForgotPasswordQueryResponse, error, struct{}]](c, &query)
if res.IsErr() {
return res.GetErr()
}
return c.JSON(http.StatusOK, res)
}