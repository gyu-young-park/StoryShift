package markdown

import (
	"path/filepath"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gyu-young-park/StoryShift/pkg/file"
)

const (
	MARKDONW_CONTENT_WITH_IMAGE = `c또는 Rust source code는 eBPF bytecode로 컴파일된다. 이 eBPF bytecode는 JIT compile되거나 interpreted되어 native machine code 명령어로 변환된다. 다음의 그림을 참고하자.  
![](https://velog.velcdn.com/images/chappi/post/42f30731-7b86-4713-a429-3acc63d288a1/image.png)

eBPF program은 eBPF bytecode 명령어 셋으로 assembly로 programming할 수 있지만, 사람이 읽을 수 있고, programming하기 좋은 c나 rust와 같은 언어로 먼저 작성하고 bytecode로 만들어 실행하는 것이 좋다. `
)

func TestMarkdownImageeMatcher(t *testing.T) {
	imageHandler := NewMarkdownImageHandler()
	pictures := imageHandler.GetImageList(MARKDONW_CONTENT_WITH_IMAGE)
	assert.Equal(t, "https://velog.velcdn.com/images/chappi/post/42f30731-7b86-4713-a429-3acc63d288a1/image.png", pictures[0])
}

func TestReplaceAllMarkdownImageUrl(t *testing.T) {
	imageHandler := NewMarkdownImageHandler()
	replacedMarkdownContent := imageHandler.ReplaceAllImageUrlOfContensWithPrefix("image-prefix", MARKDONW_CONTENT_WITH_IMAGE)
	pictures := imageHandler.GetImageList(replacedMarkdownContent)
	assert.Equal(t, "image-prefix-0.png", pictures[0])
}

func TestDownloadImageWithUrl(t *testing.T) {
	imageHandler := NewMarkdownImageHandler()
	images := imageHandler.GetImageList(MARKDONW_CONTENT_WITH_IMAGE)
	fh := file.NewFileHandler()
	defer fh.Close()

	downloadImageWithUrlRespModel, err := imageHandler.DownloadImageWithUrl(fh, []DownloadImageWithUrlReqModel{
		DownloadImageWithUrlReqModel{
			Url:           images[0],
			ImageFileName: "ebpf",
		},
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, 0, len(downloadImageWithUrlRespModel.FailedToDownloadImageUrlList))
	assert.Equal(t, 1, len(downloadImageWithUrlRespModel.ImageFilePathList))
	assert.Equal(t, "ebpf.png", filepath.Base(downloadImageWithUrlRespModel.ImageFilePathList[0]))
}

func TestMarkdownImageProcessBestScenario(t *testing.T) {
	imageHandler := NewMarkdownImageHandler()
	imageList := imageHandler.GetImageList(MARKDONW_CONTENT_WITH_IMAGE)

	fh := file.NewFileHandler()
	defer fh.Close()

	downloadImageWithUrlRespModel, err := imageHandler.DownloadImageWithUrl(fh, []DownloadImageWithUrlReqModel{
		DownloadImageWithUrlReqModel{
			Url:           imageList[0],
			ImageFileName: "ebpf-replace",
		},
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, 0, len(downloadImageWithUrlRespModel.FailedToDownloadImageUrlList))
	assert.Equal(t, 1, len(downloadImageWithUrlRespModel.ImageFilePathList))
	assert.Equal(t, "ebpf-replace.png", filepath.Base(downloadImageWithUrlRespModel.ImageFilePathList[0]))
}
