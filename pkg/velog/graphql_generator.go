package velog

import (
	"encoding/json"
	"fmt"
)

var graphQLQuery = graphQLQueryManager{}

type graphQLQueryManager struct {
}

func (qm graphQLQueryManager) readPost(username, urlSlug string) string {
	q := `
	query ReadPost($username: String, $url_slug: String) {
	  post(username: $username, url_slug: $url_slug) {
		id
		title
		released_at
		updated_at
		tags
		body
		short_description
		is_markdown
		is_private
		is_temp
		thumbnail
		comments_count
		url_slug
		likes
		liked
		user {
		  id
		  username
		  is_followed
		  profile {
			id
			display_name
			thumbnail
			short_bio
			profile_links
		  }
		  velog_config {
			title
		  }
		}
		comments {
		  id
		  user {
			id
			username
			profile {
			  id
			  thumbnail
			  display_name
			}
		  }
		  text
		  replies_count
		  level
		  created_at
		  deleted
		}
		series {
		  id
		  name
		  url_slug
		  series_posts {
			id
			post {
			  id
			  title
			  url_slug
			  user {
				id
				username
			  }
			}
		  }
		}
		linked_posts {
		  previous {
			id
			title
			url_slug
			user {
			  id
			  username
			}
		  }
		  next {
			id
			title
			url_slug
			user {
			  id
			  username
			}
		  }
		}
	  }
	}
	`
	body := graphQLRequestBody{
		OperationName: "ReadPost",
		Variables: variables{
			Username: username,
			UrlSlug:  urlSlug,
		},
		Query: q,
	}

	rawGraphQL, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(rawGraphQL)
}
