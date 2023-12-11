// Register godoc
// @Security BearerAuth
// @Summary      Register
// @Description  Register
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        loginReq body queries.RegisterQuery true "Username"
// @Success      200  {object}   baseres.SwaggerSuccessRes[queries.RegisterQueryResponse] "OK. On success."
// @Failure      400  {object}   baseres.SwaggerValidationErrRes "Bad Request. On any validation error."
// @Failure      401  {object}   baseres.SwaggerUnauthorizedErrRes "Unauthorized."
// @Failure      500  {object}   baseres.SwaggerInternalErrRes "Internal Server Error."
// @Router       /v1/auth/login [post]
func (ac *AuthController) Register(c echo.Context) error {
var query queries.RegisterQuery
if err := utils.BindAndValidate(c, &query); err != nil {
return err
}
res := mediatr.Send[*queries.RegisterQuery, *baseres.Result[*queries.RegisterQueryResponse, error, struct{}]](c, &query)
if res.IsErr() {
return res.GetErr()
}
return c.JSON(http.StatusOK, res)
}