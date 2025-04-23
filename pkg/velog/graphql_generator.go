package velog

import (
	"encoding/json"
	"fmt"
)

func getLastPostHistoryRawGraphQL(postId string) string {
	query := `
		query GetLastPostHistory($post_id: ID!) {
			lastPostHistory(post_id: $post_id) {
				id
				title
				body
				created_at
				is_markdown
				__typename
			}
		}`
	body := graphQLRequestBody{
		OperationName: "GetLastPostHistory",
		Variables: variables{
			PostId: postId,
		},
		Query: query,
	}

	rawGraphQL, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(rawGraphQL)
}
