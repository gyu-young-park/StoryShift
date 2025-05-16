package httpclient

import (
	"io"
	"net/http"
	"strings"
)

type GetRequestParam struct {
	URL   string
	Query map[string]string
}

func (p GetRequestParam) MakeURL() string {
	var sb strings.Builder
	sb.WriteString(p.URL)

	if p.Query != nil {
		isFirst := true
		for k, v := range p.Query {
			if isFirst {
				sb.WriteString("?")
				isFirst = false
			} else {
				sb.WriteString("&")
			}

			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(v)
		}
	}
	return sb.String()
}

type HeadRequestParam struct {
	GetRequestParam
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

func MakeResponseParam(resp http.Response) (ResponseParam, error) {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseParam{StatusCode: resp.StatusCode}, err
	}
	return ResponseParam{StatusCode: resp.StatusCode, Body: respBody}, nil
}
