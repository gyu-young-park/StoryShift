package servicevelog

import (
	"path/filepath"
	"testing"

	"github.com/gyu-young-park/StoryShift/pkg/file"
	"github.com/stretchr/testify/assert"
)

func TestSanitizeBasePathSpecialCaseWithSuccess(t *testing.T) {
	data := "hello world it's '/' test"
	expect := "hello world it's '-' test"

	sanitized, isSanitized := sanitizeBasePathSpecialCase(data)

	assert.Equal(t, expect, sanitized)
	assert.True(t, isSanitized)
}

func TestMarkdownImageeMatcher(t *testing.T) {
	data := `c또는 Rust source code는 eBPF bytecode로 컴파일된다. 이 eBPF bytecode는 JIT compile되거나 interpreted되어 native machine code 명령어로 변환된다. 다음의 그림을 참고하자.  
![](https://velog.velcdn.com/images/chappi/post/42f30731-7b86-4713-a429-3acc63d288a1/image.png)

eBPF program은 eBPF bytecode 명령어 셋으로 assembly로 programming할 수 있지만, 사람이 읽을 수 있고, programming하기 좋은 c나 rust와 같은 언어로 먼저 작성하고 bytecode로 만들어 실행하는 것이 좋다. `

	pictures := markdownImageMatcher(data)
	assert.Equal(t, "https://velog.velcdn.com/images/chappi/post/42f30731-7b86-4713-a429-3acc63d288a1/image.png", pictures[0])
}

func TestDownloadImageWithUrl(t *testing.T) {
	data := `c또는 Rust source code는 eBPF bytecode로 컴파일된다. 이 eBPF bytecode는 JIT compile되거나 interpreted되어 native machine code 명령어로 변환된다. 다음의 그림을 참고하자.  
![](https://velog.velcdn.com/images/chappi/post/42f30731-7b86-4713-a429-3acc63d288a1/image.png)

eBPF program은 eBPF bytecode 명령어 셋으로 assembly로 programming할 수 있지만, 사람이 읽을 수 있고, programming하기 좋은 c나 rust와 같은 언어로 먼저 작성하고 bytecode로 만들어 실행하는 것이 좋다. `

	images := markdownImageMatcher(data)
	imagesWithName := map[string]string{
		"ebpf": images[0],
	}
	fh := file.NewFileHandler()
	defer fh.Close()

	imageFilePathList := downloadImageWithUrl(fh, imagesWithName)
	assert.Equal(t, 1, len(imageFilePathList))
	assert.Equal(t, "ebpf.png", filepath.Base(imageFilePathList[0]))
}
