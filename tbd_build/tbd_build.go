package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/JimGaylard/tbd/tbd_build/build"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println("Builds the tbd project")
		fmt.Printf("%s [options...] <tree>", os.Args[0])
		flag.PrintDefaults()
	}

	config.ciDir = flag.String("ci-dir", "ci", "ci directory")
}

type Config struct {
	ciDir *string
}

var config Config

func main() {
	flag.Parse()

	if len(os.Args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	err := build.Build(os.Args[1], *config.ciDir)
	if err != nil {
		log.Fatal(err)
	}
	defer rmBuildDir(buildDir)
func rmBuildDir(buildDir string) {
	if err := os.RemoveAll(buildDir); err != nil {
		log.Fatal(err)
	}
}
