package markdown

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/gyu-young-park/StoryShift/internal/httpclient"
	"github.com/gyu-young-park/StoryShift/pkg/file"
	"github.com/gyu-young-park/StoryShift/pkg/log"
)

const (
	IMAGE_REGEX_MATCHER = `!\[[^\]]*\]\(([^)]+)\)`
	IMAGE_DIR_KEY       = "IMAGE_DIRECTORY_PATH_KEY"
)

type MarkdownImageHandlable interface {
	GetImageList(contents string) []string
	ReplaceAllImageUrlOfContensWithPrefix(imageNamePrefix string, contents string) string
	DownloadImageWithUrl(fh *file.FileHandler, imageUrls map[string]string) ([]string, []string, error)
}

type MarkdownImageManipulator struct {
	lock               *sync.Mutex
	handler            MarkdownImageHandlable
	imageNameAndURLMap map[string]string
	fileHandler        *file.FileHandler
}

func NewMarkdownImageManipulator(handler MarkdownImageHandlable) *MarkdownImageManipulator {
	return &MarkdownImageManipulator{
		lock:               &sync.Mutex{},
		handler:            handler,
		imageNameAndURLMap: make(map[string]string),
		fileHandler:        file.NewFileHandler(),
	}
}

func (m *MarkdownImageManipulator) Replace(title string, contents string) string {
	m.lock.Lock()
	defer m.lock.Unlock()
	imageTagList := m.handler.GetImageList(contents)
	if len(imageTagList) == 0 {
		return contents
	}

	encodedTitle := base64.RawURLEncoding.EncodeToString([]byte(title))
	for i, image := range imageTagList {
		m.imageNameAndURLMap[fmt.Sprintf("%s-%v", encodedTitle, i)] = image
	}

	return m.handler.ReplaceAllImageUrlOfContensWithPrefix(fmt.Sprintf("%s/%s", IMAGE_DIR_KEY, encodedTitle), contents)
}

func (m *MarkdownImageManipulator) DownloadAsZip(username string) (*os.File, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	logger := log.GetLogger()

	imageFileList := []*os.File{}
	imageFilePath, failedToDownloadTitleList, err := m.handler.DownloadImageWithUrl(m.fileHandler, m.imageNameAndURLMap)
	if err != nil {
		logger.Errorf("failed to download %s", err.Error())
		return nil, err
	}
	fmt.Println(failedToDownloadTitleList)

	fmt.Println()

	for _, image := range imageFilePath {
		if image != "" {
			imageFileList = append(imageFileList, m.fileHandler.GetFileWithLocked(image))
		}
	}

	imageZipname, err := m.fileHandler.CreateZipFile(file.ZipFile{
		FileMeta: file.FileMeta{
			Name:      fmt.Sprintf("%s-%s", username, "image"),
			Extention: "zip",
		},
		Files: imageFileList,
	})

	if err != nil {
		return nil, err
	}

	imageZipfile := m.fileHandler.GetFileWithLocked(imageZipname)
	imageZipfile.Seek(0, 0)

	return imageZipfile, nil
}

func (m *MarkdownImageManipulator) Done() {
	m.imageNameAndURLMap = map[string]string{}
	m.fileHandler.Close()
	m.fileHandler = nil
}

type MarkdownImageHandler struct {
	matcher                   *regexp.Regexp
	failedToDownloadTitleList []string
}

func NewMarkdownImageHandler() *MarkdownImageHandler {
	return &MarkdownImageHandler{
		matcher: regexp.MustCompile(IMAGE_REGEX_MATCHER),
	}
}

func (m *MarkdownImageHandler) GetImageList(contents string) []string {
	matches := m.matcher.FindAllStringSubmatch(contents, -1)

	images := []string{}
	for _, match := range matches {
		images = append(images, match[1])
	}

	return images
}

func (m *MarkdownImageHandler) ReplaceAllImageUrlOfContensWithPrefix(imageNamePrefix string, contents string) string {
	logger := log.GetLogger()
	index := 0
	return m.matcher.ReplaceAllStringFunc(contents, func(match string) string {
		imageName := fmt.Sprintf("%v-%v", imageNamePrefix, index)
		index++

		submatch := m.matcher.FindStringSubmatch(match)
		urlIndex := len(submatch) - 1
		logger.Debugf("change image url from [%s] to [%s]", submatch[urlIndex-1], imageName)
		alt := match[:len(match)-len(submatch[urlIndex])-1]
		return alt + imageName + filepath.Ext(match)
	})
}

func (m *MarkdownImageHandler) DownloadImageWithUrl(fh *file.FileHandler, imageUrls map[string]string) ([]string, []string, error) {
	logger := log.GetLogger()
	if imageUrls == nil {
		return []string{}, []string{}, fmt.Errorf("there is no imageUrls map")
	}

	files := []file.File{}
	for replaceName, imageUrl := range imageUrls {
		resp, err := httpclient.Get(httpclient.GetRequestParam{
			URL: imageUrl,
		})
		if err != nil {
			m.failedToDownloadTitleList = append(m.failedToDownloadTitleList, imageUrl)
		}

		_, ext := file.SplitFilenameWithNameAndExt(imageUrl)
		files = append(files, file.File{
			FileMeta: file.FileMeta{
				Name:      replaceName,
				Extention: ext,
			},
			Content: string(resp.Body),
		})
	}

	imageFileList := []string{}
	for _, file := range files {
		imageFilePath, err := fh.CreateFile(file)
		if err != nil {
			return []string{}, m.failedToDownloadTitleList, fmt.Errorf("failed to create image file: %s", file.GetFilename())
		}
		imageFileList = append(imageFileList, imageFilePath)
		logger.Infof("image file created: %s", imageFilePath)
	}

	return imageFileList, m.failedToDownloadTitleList, nil
}

type DefaultMarkdownImageHandler struct {
}

func NewDefaultMarkdownImageHandler() *DefaultMarkdownImageHandler {
	return &DefaultMarkdownImageHandler{}
}

func (m *DefaultMarkdownImageHandler) GetImageList(contents string) []string {
	return []string{}
}

func (m *DefaultMarkdownImageHandler) ReplaceAllImageUrlOfContensWithPrefix(imageNamePrefix string, contents string) string {
	return ""
}

func (m *DefaultMarkdownImageHandler) DownloadImageWithUrl(fh *file.FileHandler, imageUrls map[string]string) ([]string, error) {
	return []string{}, nil
}
