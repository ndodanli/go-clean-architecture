package pg

type PaginationQuery struct {
	Page       int    `query:"page" validate:"required"`
	PageSize   int    `query:"pageSize" validate:"required"`
	OrderBy    string `query:"orderBy"`
	SortBy     string `query:"sortBy"`
	SearchTerm string `query:"searchTerm"`
}
