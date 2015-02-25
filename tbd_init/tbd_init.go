package main

import (
	"fmt"

	initRepo "github.com/JimGaylard/tbd/tbd_init/init"
)

func main() {
	if err := initRepo.Init(); err != nil {
		panic(err)
	}
	fmt.Println("success")
}
