package httpclient

import "io"

type PostRequestParam struct {
	URL         string
	Body        io.Reader
	ContentType string
}
