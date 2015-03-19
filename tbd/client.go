package main

import (
	"fmt"
	"os"

	"github.com/tbd-ci/tbd/builder"
)

func main() {
	if os.Args[1] == "build" {
		if err := builder.Build(); err != nil {
			fmt.Println(err)
		}
	}
}
