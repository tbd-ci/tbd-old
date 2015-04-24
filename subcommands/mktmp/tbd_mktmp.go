package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tbd-ci/tbd/subcommands/mktmp/mktmp"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println("Takes a git tree object and checks it out to a new tmp dir.")
		fmt.Printf("%s [options...] <tree>", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	if len(os.Args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	tree := os.Args[1]

	tmpDir, err := mktmp.CheckoutTmp(tree)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tmpDir)
}
