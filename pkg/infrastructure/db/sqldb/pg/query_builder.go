package pg

import (
	"fmt"
	"strings"
)

type QueryString struct {
	string               string
	firstWhereCall       bool
	firstSetCall         bool
	index                int
	pgxArgs              []interface{}
	skipDeleted          bool
	whereGroupIndex      int
	groupIndexHasBeenSet bool
	groupWhereOpen       bool
}

func NewQueryString(query string) *QueryString {
	return &QueryString{
		string:               query,
		index:                1,
		pgxArgs:              []interface{}{},
		skipDeleted:          true,
		whereGroupIndex:      0,
		groupIndexHasBeenSet: false,
		groupWhereOpen:       false,
		firstWhereCall:       true,
		firstSetCall:         true,
	}
}

func (q *QueryString) AddSet(operation string, conditionName string, arg interface{}) *QueryString {
	if q.firstSetCall {
		q.string = fmt.Sprintf("%s SET %s = $%d", q.string, conditionName, q.index)
		q.firstSetCall = false
	} else {
		q.string = fmt.Sprintf("%s %s %s = $%d", q.string, operation, conditionName, q.index)
	}
	q.pgxArgs = append(q.pgxArgs, arg)
	q.index++

	return q
}

func (q *QueryString) AddWhere(operation string, conditionName string, operator string, arg interface{}) *QueryString {
	if q.firstWhereCall && q.skipDeleted {
		q.string = fmt.Sprintf("%s WHERE deleted_at = '0001-01-01T00:00:00Z' AND %s %s $%d", q.string, conditionName, operator, q.index)
		q.firstWhereCall = false
	} else {
		q.string = fmt.Sprintf("%s %s %s = $%d", q.string, operation, conditionName, q.index)
	}
	q.pgxArgs = append(q.pgxArgs, arg)
	q.index++

	return q
}

func (q *QueryString) StartGroupedWhere(operation string) *QueryString {
	if q.firstWhereCall && q.skipDeleted {
		q.string = fmt.Sprintf("%s WHERE deleted_at = '0001-01-01T00:00:00Z' AND (", q.string)
		q.firstWhereCall = false
	} else {

		q.string = fmt.Sprintf("%s %s (", q.string, operation)
	}

	q.groupWhereOpen = true
	return q
}

func (q *QueryString) AddToGroupedWhere(operation string, condition string, operator string, arg interface{}, groupIndex int) *QueryString {
	isGroupIndexProvided := groupIndex != 0
	if !isGroupIndexProvided {
		groupIndex = q.index
	} else if arg != nil {
		panic("arg must be nil if groupIndex is not 0")
	}

	if !isGroupIndexProvided && arg == nil {
		panic("arg cannot be nil if groupIndex is 0")
	}

	if !q.groupIndexHasBeenSet {
		operation = ""
	}

	q.string = fmt.Sprintf("%s %s %s %s $%d", q.string, operation, condition, operator, groupIndex)
	if !isGroupIndexProvided || arg != nil {
		q.pgxArgs = append(q.pgxArgs, arg)
	}
	if !isGroupIndexProvided {
		q.index++
	}

	if !q.groupIndexHasBeenSet {
		q.whereGroupIndex = groupIndex
		q.groupIndexHasBeenSet = true
	}

	return q
}

func (q *QueryString) CloseGroupedWhere() *QueryString {
	if !q.groupWhereOpen {
		panic("grouped where is not open")
	}
	q.string = fmt.Sprintf("%s)", q.string)
	q.groupIndexHasBeenSet = false
	q.groupWhereOpen = false
	return q

}

func (q *QueryString) Paginate(pq *PaginationQuery, ordering bool) *QueryString {
	orderByClause := ""
	if ordering {
		if pq.OrderBy == "" {
			pq.OrderBy = "created_at"
		}

		if pq.SortBy == "" {
			pq.SortBy = "DESC"
		}

		orderByClause = fmt.Sprintf(" ORDER BY %s %s", pq.OrderBy, pq.SortBy)
	}

	offsetLimitClause := fmt.Sprintf(" OFFSET %d LIMIT %d", (pq.Page-1)*pq.PageSize, pq.PageSize)
	q.string = strings.TrimSpace(fmt.Sprintf("%s%s%s", q, orderByClause, offsetLimitClause))
	return q
}

func (q *QueryString) SkipDeleted() *QueryString {
	// Define the condition for deleted_at
	deletedAtCondition := "deleted_at = '0001-01-01T00:00:00Z'"

	// Split the query into parts
	parts := strings.SplitN(q.string, "WHERE", 2)

	// Construct the final query
	var fullQuery string
	if len(parts) > 1 {
		// If WHERE clause exists, insert the condition after WHERE
		fullQuery = fmt.Sprintf("%s WHERE %s AND %s", parts[0], deletedAtCondition, parts[1])
	} else {
		// If no WHERE clause exists, add the WHERE clause with the condition after table name
		fullQuery = fmt.Sprintf("%s WHERE %s", parts[0], deletedAtCondition)
	}
	q.string = strings.TrimSpace(fullQuery)
	return q
}

func (q *QueryString) CurrentIndex() int {
	return q.index
}

func (q *QueryString) CurrentWhereGroupIndex() int {
	return q.whereGroupIndex
}

func (q *QueryString) String() string {
	return q.string
}

func (q *QueryString) Args() []interface{} {
	return q.pgxArgs
}

func (q *QueryString) Copy() *QueryString {
	return &QueryString{
		string:               q.string,
		index:                q.index,
		pgxArgs:              q.pgxArgs,
		skipDeleted:          q.skipDeleted,
		whereGroupIndex:      q.whereGroupIndex,
		groupIndexHasBeenSet: q.groupIndexHasBeenSet,
		groupWhereOpen:       q.groupWhereOpen,
	}
}

func (q *QueryString) GetCountQuery() string {
	fromIndex := strings.Index(q.string, "FROM")
	return fmt.Sprintf("SELECT COUNT(*) %s", q.string[fromIndex:len(q.string)])
}
