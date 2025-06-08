package servicevelog

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gyu-young-park/StoryShift/internal/httpclient"
	"github.com/gyu-young-park/StoryShift/pkg/file"
	"github.com/gyu-young-park/StoryShift/pkg/log"
)

func sanitizeBasePathSpecialCase(filename string) (string, bool) {
	re := regexp.MustCompile(`[/]`)
	matched := re.MatchString(filename)
	sanitize := re.ReplaceAllString(filename, "-")
	return sanitize, matched
}

func markdownImageMatcher(contents string) []string {
	re := regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(contents, -1)

	images := []string{}
	for _, match := range matches {
		images = append(images, match[1])
	}

	return images
}

func downloadImageWithUrl(fh *file.FileHandler, imageUrls map[string]string) []string {
	logger := log.GetLogger()

	if imageUrls == nil {
		logger.Error("there is no imageUrls map")
		return []string{}
	}

	imageFileList := []string{}
	for imageName, imageUrl := range imageUrls {
		resp, err := httpclient.Get(httpclient.GetRequestParam{
			URL: imageUrl,
		})

		if err != nil {
			logger.Errorf("failed to download image: %s", imageUrl)
			continue
		}

		ext := strings.TrimPrefix(filepath.Ext(imageUrl), ".")
		imageFilePath, err := fh.CreateFile(file.File{
			FileMeta: file.FileMeta{
				Name:      imageName,
				Extention: ext,
			},
			Content: string(resp.Body),
		})

		if err != nil {
			logger.Errorf("failed to create image file: %s", imageUrl)
			continue
		}
		imageFileList = append(imageFileList, imageFilePath)
	}

	return imageFileList
}
