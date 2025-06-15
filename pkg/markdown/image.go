package markdown

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/gyu-young-park/StoryShift/pkg/file"
	"github.com/gyu-young-park/StoryShift/pkg/log"
)

const (
	IMAGE_REGEX_MATCHER = `!\[[^\]]*\]\(([^)]+)\)`
	IMAGE_DIR_KEY       = "IMAGE_DIRECTORY_PATH_KEY"
)

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
	imageFilePath, failedImageList, err := m.handler.DownloadImageWithUrl(m.fileHandler, m.imageNameAndURLMap)
	if err != nil {
		logger.Errorf("failed to download %s", err.Error())
		return nil, err
	}

	if len(failedImageList) > 0 {
		bFailedImageList, _ := json.Marshal(failedImageList)
		downloadFailListFilePath, err := m.fileHandler.CreateFile(file.File{
			FileMeta: file.FileMeta{
				Name:      "download_fail_list",
				Extention: "json",
			},
			Content: string(bFailedImageList),
		})

		if err != nil {
			logger.Errorf("failed to create download fail list, err: %s", err.Error())
		} else {
			imageFileList = append(imageFileList, m.fileHandler.GetFileWithLocked(downloadFailListFilePath))
		}
	}

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
