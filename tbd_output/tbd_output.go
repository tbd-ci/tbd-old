package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/JimGaylard/tbd/tbd_output/output"
)

type Config struct {
	stream *string
}

var config Config

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println("Returns output for a given treeish git object")
		fmt.Printf("%s [options...] <treeish>", os.Args[0])
		flag.PrintDefaults()
	}

	config.stream = flag.String("only", "combined", "only display stdout|stderr or combined")
	config.build = flag.String("prompt-for-build", "latest", "display output of which build number/run")
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println(flag.Args())
	}

	if err := output.Display(flag.Args()[0], &config); err != nil {
		panic(err)
	}
}
