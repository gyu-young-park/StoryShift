package httpclient

import "io"

type GetRequestParam struct {
	URL   string
	Query map[string]string
}

type PostRequestParam struct {
	URL         string
	Body        io.Reader
	ContentType string
}

type ResponseParam struct {
	StatusCode int
	Body       []byte
}
