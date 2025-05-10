package httpclient

import (
	"io"
	"net/http"
	"strings"
	"time"
)

var defaultHTTPClient = NewDefaultHttpClient(time.Second*10, newDefaultTransport())

type defaultHttpClient struct {
	*http.Client
}

func NewDefaultHttpClient(timeout time.Duration, customTransport http.RoundTripper) *defaultHttpClient {
	return &defaultHttpClient{
		Client: &http.Client{
			Timeout:   timeout,
			Transport: customTransport,
		},
	}
}

func Get(param GetRequestParam) (ResponseParam, error) {
	var sb strings.Builder
	sb.WriteString(param.URL)

	isFirst := true
	for k, v := range param.Query {
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

	url := sb.String()
	resp, err := defaultHTTPClient.Get(url)
	if err != nil {
		return ResponseParam{StatusCode: resp.StatusCode}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseParam{StatusCode: resp.StatusCode}, err
	}
	return ResponseParam{StatusCode: resp.StatusCode, Body: respBody}, nil

}

func Post(param PostRequestParam) (ResponseParam, error) {
	resp, err := defaultHTTPClient.Post(param.URL, param.ContentType, param.Body)
	if err != nil {
		return ResponseParam{StatusCode: resp.StatusCode}, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseParam{StatusCode: resp.StatusCode}, err
	}
	return ResponseParam{StatusCode: resp.StatusCode, Body: respBody}, nil
}
