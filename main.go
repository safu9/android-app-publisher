package main

import (
	"fmt"
	"os"

	"github.com/safu9/android-app-publisher/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
