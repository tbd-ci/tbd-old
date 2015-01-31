package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	path, err := mktmp_d()
	if err != nil {
		fmt.Println(err)
		log.Fatal("mktmp_d: %s", err)
	}

	tmpDir, err := checkout_tmp(path)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("checkout: %s", err)
	}
	fmt.Print(tmpDir)
}

func mktmp_d() (string, error) {
	out, err := exec.Command("mktemp", "-d").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func checkout_tmp(tree string) (string, error) {
	tmp_path, err := mktmp_d()
	if err != nil {
		return "", err
	}

	cmd1 := exec.Command("git", "archive", tree)
	cmd2 := exec.Command("tar", "-x", "-C", tmp_path)

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

	return tmp_path, err
}
