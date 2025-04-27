package httpclient

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gyu-young-park/VelogStoryShift/pkg/log"
)

func newDefaultTransport() *defaultTransport {
	dialer := &net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 5 * time.Second,
	}

	return &defaultTransport{
		core: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          50,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   2 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

type defaultTransport struct {
	core http.RoundTripper
}

func (d *defaultTransport) roundTripWithRetry(req *http.Request, maxRetries int) (*http.Response, error) {
	logger := log.GetLogger()
	for i := 0; i < maxRetries; i++ {
		res, err := d.core.RoundTrip(req)
		if err == nil {
			return res, nil
		}

		logger.Debugf("try: %v, retry request URL: %s", i+1, req.URL)
	}
	return nil, fmt.Errorf("failed to request URL: %s", req.URL)
}

func (d *defaultTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	logger := log.GetLogger()

	now := time.Now()
	logger.Debugf("Start to request URL: %s\n", req.URL)
	res, err := d.roundTripWithRetry(req, 3)
	logger.Debugf("HTTP Client URL: %s, latency: %v s\n", req.URL, time.Since(now).Seconds())
	if err != nil {
		return nil, err
	}
	return res, nil
}
