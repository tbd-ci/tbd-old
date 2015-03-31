package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tbd-ci/tbd/tbd_mktmp/mktmp"
	"github.com/tbd-ci/tbd/tbd_note/note"
)

var cmdVal int

type BuildConfig struct {
	Treeish  string
	CiDir    string
	BuildDir string
}

func rmTmp(tmpDir string) error {
	if err := os.RemoveAll(tmpDir); err != nil {
		return err
	}
	return nil
}

func Build(config BuildConfig) error {
	walkFunc := func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		// Only run for each of "ci/stage" subdirectories
		if !info.IsDir() || filepath.Base(path) == config.CiDir {
			return nil
		}

		tmpDir, err := mktmp.CheckoutTmp(config.Treeish)
		if err != nil {
			return err
		}
		defer rmTmp(tmpDir)

		// Get the relative path of each command to run as we will be running
		// it in its own tmp directory (not in `path`)
		ciCmd := filepath.Join(path, "run")
		relCmd, err := filepath.Rel(config.BuildDir, ciCmd)
		if err != nil {
			return err
		}

		refPrefix := filepath.Join("refs", "builds", config.Treeish, filepath.Base(path))
		cmd := exec.Command("tbd_run", "-prefix", refPrefix, filepath.Join(tmpDir, relCmd))

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return err
		}

		return nil
	}

	err := filepath.Walk(filepath.Join(config.BuildDir, config.CiDir), walkFunc)
	if err != nil {
		return err
	}

	success := fmt.Sprintln("Build was successful")
	failure := fmt.Sprintln("Build failed")

	if cmdVal == 0 {
		note.WriteNote(
			config.Treeish,
			success,
		)
	} else {
		note.WriteNote(
			config.Treeish,
			failure,
		)
	}

	return nil
}
