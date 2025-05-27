package servicestatus

import (
	"context"
	"time"

	"github.com/gyu-young-park/StoryShift/internal/config"
	"github.com/gyu-young-park/StoryShift/internal/httpclient"
	"github.com/gyu-young-park/StoryShift/pkg/log"
)

func (s *StatusService) Ready() bool {
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

	if !s.redisCheck() {
		return false
	}

	return true
}

func (s *StatusService) redisCheck() bool {
	logger := log.GetLogger()
	ctx, close := context.WithTimeout(context.Background(), time.Second*5)
	defer close()

	_, err := s.redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Error("unstable redis server: " + err.Error())
		return false
	}

	return true
}
