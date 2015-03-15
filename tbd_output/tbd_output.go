package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/JimGaylard/tbd/tbd_output/output"
)

var config output.Config

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println("Returns output for a given treeish git object")
		fmt.Printf("%s [options...] <treeish>", os.Args[0])
		flag.PrintDefaults()
	}

	config.Stream = *flag.String("only", "", "only display stdout|stderr or combined")
	config.Build = *flag.String("build", "latest", "display output of which build number/run")
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	if err := output.Display(flag.Args()[0], config); err != nil {
		panic(err)
	}
}
