package mktmp

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func CheckoutTmp(tree string) (string, error) {
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
