package httpclient

import (
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

func Head(param HeadRequestParam) (ResponseParam, error) {
	resp, err := defaultHTTPClient.Head(param.MakeURL())
	if err != nil {
		return ResponseParam{}, err
	}

	defer resp.Body.Close()

	return MakeResponseParam(*resp)
}

func Get(param GetRequestParam) (ResponseParam, error) {
	resp, err := defaultHTTPClient.Get(param.MakeURL())
	if err != nil {
		return ResponseParam{}, err
	}

	defer resp.Body.Close()

	return MakeResponseParam(*resp)
}

func Post(param PostRequestParam) (ResponseParam, error) {
	resp, err := defaultHTTPClient.Post(param.URL, param.ContentType, param.Body)
	if err != nil {
		return ResponseParam{}, err
	}

	defer resp.Body.Close()

	return MakeResponseParam(*resp)
}
