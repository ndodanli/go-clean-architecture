package appuserrepo

type GetOnlyIdRepoRes struct {
	ID       int64  `json:"id"`
	Password string `json:"password"`
}
