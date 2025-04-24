package velog

type velogOperationEnum string

const (
	VELOG_OPERATION_POSTS     velogOperationEnum = "Posts"
	VELOG_OPERATION_READ_POST velogOperationEnum = "ReadPost"
)

type velogQueryEnum string

const (
	VELOG_QUERY_READ_POST velogQueryEnum = `query ReadPost($username: String, $url_slug: String) {
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
	}`

	VELOG_QUERY_POSTS velogQueryEnum = `query Posts($cursor: ID, $username: String, $temp_only: Boolean, $tag: String, $limit: Int) {
		posts(cursor: $cursor, username: $username, temp_only: $temp_only, tag: $tag, limit: $limit) {
			id
			title
			short_description
			thumbnail
			user {
			id
			username
			profile {
				id
				thumbnail
				__typename
			}
			__typename
			}
			url_slug
			released_at
			updated_at
			comments_count
			tags
			is_private
			likes
			__typename
		}
	}`
)
