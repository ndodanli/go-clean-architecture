package queries

import "context"

type TestQueryHandler struct {
}

type TestQuery struct {
	TestID string
}

func NewTestQueryHandler() *TestQueryHandler {
	return &TestQueryHandler{}
}

type TestQueryResponse struct {
	TestIDRes string `json:"testIDRes"`
}

func (c *TestQueryHandler) Handle(ctx context.Context, query *TestQuery) (*TestQueryResponse, error) {

	return &TestQueryResponse{
		TestIDRes: "dsadsa",
	}, nil
}
