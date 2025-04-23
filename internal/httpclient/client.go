package httpclient

import (
	"io"
	"net/http"
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

func Post(param PostRequestParam) ([]byte, error) {
	resp, err := defaultHTTPClient.Post(param.URL, param.ContentType, param.Body)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return respBody, nil
}
