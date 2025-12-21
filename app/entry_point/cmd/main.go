package main

import (
	"fmt"
	"os"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/cmd"
)

func main() {
	cmd := cmd.NewCmd()
	cmd.Init()

	if err := cmd.Cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
