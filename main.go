package main

import (
	"fmt"
	"os"

	"github.com/trapacska/repotool/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
