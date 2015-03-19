package main

import (
	initRepo "github.com/tbd-ci/tbd/tbd_init/init"
)

func main() {
	if err := initRepo.Init(); err != nil {
		panic(err)
	}
}
