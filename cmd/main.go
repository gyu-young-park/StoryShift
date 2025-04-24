package main

import (
	"fmt"
	"os"

	"github.com/gyu-young-park/VelogStoryShift/internal/config"
	"github.com/gyu-young-park/VelogStoryShift/pkg/velog"
)

func main() {
	velogAPI := velog.NewVelogAPI(config.Manager.VelogConfig.URL, "chappi")
	resp, err := velogAPI.GetPost("ElasticSearch-정리-9일차-Aggregation")
	if err != nil {
		fmt.Printf("failed to get post content of velog: %s", err)
		return
	}
	fmt.Println(resp.Body)
	// inputFile := "data.txt"
	outputFile := "output.md"

	// // 파일 읽기
	// data, err := os.ReadFile(inputFile)
	// if err != nil {
	// 	fmt.Println("파일 읽기 실패:", err)
	// 	return
	// }

	// content := string(data)

	// 이스케이프된 줄바꿈 문자열("\n")을 실제 줄바꿈으로 변환
	// converted, err := strconv.Unquote(resp.Data.LastPostHistory.Body)
	// if err != nil {
	// 	fmt.Println("unquote 실패", err)
	// 	return
	// }
	// 변환된 내용을 출력 파일에 저장
	err = os.WriteFile(outputFile, []byte(resp.Body), 0644)
	if err != nil {
		fmt.Println("파일 쓰기 실패:", err)
		return
	}

	// fmt.Println("변환 완료! 결과는", outputFile)
}
