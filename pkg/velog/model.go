package velog

import "time"

type variables struct {
	UrlSlug  string `json:"url_slug"`
	PostId   string `json:"post_id"`
	Username string `json:"username"`
}

type graphQLRequestBody struct {
	OperationName string    `json:"operationName"`
	Variables     variables `json:"variables"`
	Query         string    `json:"query"`
}

type readPostModel struct {
	Data struct {
		Post struct {
			ID               string      `json:"id"`
			Title            string      `json:"title"`
			ReleasedAt       time.Time   `json:"released_at"`
			UpdatedAt        time.Time   `json:"updated_at"`
			Tags             []string    `json:"tags"`
			Body             string      `json:"body"`
			ShortDescription string      `json:"short_description"`
			IsMarkdown       bool        `json:"is_markdown"`
			IsPrivate        bool        `json:"is_private"`
			IsTemp           bool        `json:"is_temp"`
			Thumbnail        interface{} `json:"thumbnail"`
			CommentsCount    int         `json:"comments_count"`
			URLSlug          string      `json:"url_slug"`
			Likes            int         `json:"likes"`
			Liked            bool        `json:"liked"`
			User             struct {
				ID         string `json:"id"`
				Username   string `json:"username"`
				IsFollowed bool   `json:"is_followed"`
				Profile    struct {
					ID           string `json:"id"`
					DisplayName  string `json:"display_name"`
					Thumbnail    string `json:"thumbnail"`
					ShortBio     string `json:"short_bio"`
					ProfileLinks struct {
						URL      string `json:"url"`
						Email    string `json:"email"`
						Github   string `json:"github"`
						Twitter  string `json:"twitter"`
						Facebook string `json:"facebook"`
					} `json:"profile_links"`
					Typename string `json:"__typename"`
				} `json:"profile"`
				VelogConfig struct {
					Title    interface{} `json:"title"`
					Typename string      `json:"__typename"`
				} `json:"velog_config"`
				Typename string `json:"__typename"`
			} `json:"user"`
			Comments []interface{} `json:"comments"`
			Series   struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				URLSlug     string `json:"url_slug"`
				SeriesPosts []struct {
					ID   string `json:"id"`
					Post struct {
						ID      string `json:"id"`
						Title   string `json:"title"`
						URLSlug string `json:"url_slug"`
						User    struct {
							ID       string `json:"id"`
							Username string `json:"username"`
							Typename string `json:"__typename"`
						} `json:"user"`
						Typename string `json:"__typename"`
					} `json:"post"`
					Typename string `json:"__typename"`
				} `json:"series_posts"`
				Typename string `json:"__typename"`
			} `json:"series"`
			LinkedPosts struct {
				Previous struct {
					ID      string `json:"id"`
					Title   string `json:"title"`
					URLSlug string `json:"url_slug"`
					User    struct {
						ID       string `json:"id"`
						Username string `json:"username"`
						Typename string `json:"__typename"`
					} `json:"user"`
					Typename string `json:"__typename"`
				} `json:"previous"`
				Next     interface{} `json:"next"`
				Typename string      `json:"__typename"`
			} `json:"linked_posts"`
			Typename string `json:"__typename"`
		} `json:"post"`
	} `json:"data"`
}

type VelogPost struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
