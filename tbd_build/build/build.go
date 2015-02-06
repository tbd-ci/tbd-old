package build

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/JimGaylard/tbd/tbd_mktmp/mktmp"
)

func Build(tree, ciPath string) error {
	// TODO: delete tmpDir after run
	walkFunc := func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if !info.IsDir() || path == "ci" {
			return nil
		}

		tmpDir, err := mktmp.CheckoutTmp(tree)
		if err != nil {
			return err
		}

		ciCmd := filepath.Join(tmpDir, path, "run")

		refPrefix := filepath.Join("refs", "builds", tree, filepath.Base(path))
		cmd := exec.Command("tbd_run", "-propagateErrors", "-prefix", refPrefix, ciCmd)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return err
		}

		return nil
	}

	err := filepath.Walk(ciPath, walkFunc)
	if err != nil {
		return err
	}

	return nil
}
