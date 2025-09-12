package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "storyshift",       // 실행 명령어 이름
		Short: "storyshift usage", // 짧은 설명
		Long:  "storyshift is the cli for downloading velog blog",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from MyApp!")
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
