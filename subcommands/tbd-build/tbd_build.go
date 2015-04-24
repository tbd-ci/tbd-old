package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tbd-ci/tbd/subcommands/tbd-build/build"
	"github.com/tbd-ci/tbd/subcommands/tbd-mktmp/mktmp"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println("Builds the tbd project")
		fmt.Printf("%s [options...] <tree>", os.Args[0])
		flag.PrintDefaults()
	}

	buildConfig.CiDir = *flag.String("ci-dir", "ci", "ci directory")
}

var buildConfig build.BuildConfig

func main() {
	flag.Parse()

	if len(os.Args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	buildDir, err := mktmp.CheckoutTmp(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer rmBuildDir(buildDir)

	buildConfig.BuildDir = buildDir
	buildConfig.Treeish = os.Args[1]

	err = build.Build(buildConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func rmBuildDir(buildDir string) {
	if err := os.RemoveAll(buildDir); err != nil {
		log.Fatal(err)
	}
}
