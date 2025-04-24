package velog

import (
	"encoding/json"
	"fmt"
)

var graphQLQuery = graphQLQueryManager{}

type graphQLQueryManager struct {
}

func makeGraphQLQuery(operation velogOperationEnum, q velogQueryEnum, v variables) string {
	body := graphQLRequestBody{
		OperationName: string(operation),
		Variables:     v,
		Query:         string(q),
	}

	r, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(r)
}

func (qm graphQLQueryManager) posts(username, cursor string, limit int) string {
	return makeGraphQLQuery(VELOG_OPERATION_POSTS, VELOG_QUERY_POSTS, variables{
		Username: username,
		Cursor:   cursor,
		Litmit:   limit,
	})
}

func (qm graphQLQueryManager) readPost(username, urlSlug string) string {
	return makeGraphQLQuery(VELOG_OPERATION_READ_POST, VELOG_QUERY_READ_POST, variables{
		Username: username,
		UrlSlug:  urlSlug,
	})
}
