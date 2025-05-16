package servicevelog

import (
	"net/http"

	"github.com/gyu-young-park/StoryShift/internal/config"
	"github.com/gyu-young-park/StoryShift/internal/httpclient"
	"github.com/gyu-young-park/StoryShift/pkg/log"
	"github.com/gyu-young-park/StoryShift/pkg/velog"
)

func GetUserProfile(username string) (velog.VelogUserProfile, error) {
	velogApi := velog.NewVelogAPI(config.Manager.VelogConfig.ApiUrl, username)
	userProfile, err := velogApi.UserProfile()
	if err != nil {
		return velog.VelogUserProfile{}, err
	}

	return userProfile, nil
}

func IsVelogUserExists(username string) bool {
	logger := log.GetLogger()

	resp, err := httpclient.Get(httpclient.GetRequestParam{
		URL: "https://velog.io/@" + username,
	})

	if err != nil {
		logger.Errorf("failed to get whether user(%s) is valid or not, err: %s", username, err.Error())
		return false
	}

	if resp.StatusCode != http.StatusOK {
		logger.Debugf("user(%s) is not the velog user, status code: %v", username, resp.StatusCode)
		return false
	}

	return true
}
