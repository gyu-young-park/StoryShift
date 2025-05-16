package servicestatus

import (
	"github.com/gyu-young-park/StoryShift/internal/config"
	"github.com/gyu-young-park/StoryShift/internal/httpclient"
	"github.com/gyu-young-park/StoryShift/pkg/log"
)

func Ready() bool {
	logger := log.GetLogger()
	res, err := httpclient.Head(httpclient.HeadRequestParam{
		GetRequestParam: httpclient.GetRequestParam{
			URL: config.Manager.VelogConfig.Url,
		},
	})

	if err != nil {
		logger.Errorf("failed to check health of velog, err: %s", err.Error())
		return false
	}

	if res.StatusCode != 200 {
		logger.Errorf("failed to check health of velog, status code: %d", res.StatusCode)
		return false
	}

	return true
}
