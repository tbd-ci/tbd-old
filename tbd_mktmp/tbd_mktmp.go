package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

	tmpDir, err := checkoutTmp(tree)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tmpDir)
}

func checkoutTmp(tree string) (string, error) {
	tmpPath, err := ioutil.TempDir("", "tbd_build")
	if err != nil {
		return "", err
	}

	cmd1 := exec.Command("git", "archive", tree)
	cmd2 := exec.Command("tar", "-x", "-C", tmpPath)

	cmd2.Stdin, err = cmd1.StdoutPipe()
	if err != nil {
		return "", err
	}

	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr

	err = cmd2.Start()
	if err != nil {
		return "", err
	}

	err = cmd1.Run()
	if err != nil {
		return "", err
	}

	err = cmd2.Wait()
	if err != nil {
		return "", err
	}

	return tmpPath, nil
}
