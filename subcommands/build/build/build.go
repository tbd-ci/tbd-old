package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tbd-ci/tbd/subcommands/mktmp/mktmp"
	"github.com/tbd-ci/tbd/subcommands/note/note"
)

var cmdVal int

func rmTmp(tmpDir string) error {
	if err := os.RemoveAll(tmpDir); err != nil {
		return err
	}
	return nil
}

func Build(tree, buildDir, ciPath string) error {
	walkFunc := func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		// Only run for each of "ci/stage" subdirectories
		if !info.IsDir() || filepath.Base(path) == "ci" {
			return nil
		}

		tmpDir, err := mktmp.CheckoutTmp(tree)
		if err != nil {
			return err
		}
		defer rmTmp(tmpDir)

		// Get the relative path of each command to run as we will be running
		// it in its own tmp directory (not in `path`)
		ciCmd := filepath.Join(path, "run")
		relCmd, err := filepath.Rel(buildDir, ciCmd)
		if err != nil {
			return err
		}

		refPrefix := filepath.Join("refs", "builds", tree, filepath.Base(path))
		cmd := exec.Command("tbd_run", "-prefix", refPrefix, filepath.Join(tmpDir, relCmd))

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return err
		}

		return nil
	}

	err := filepath.Walk(filepath.Join(buildDir, ciPath), walkFunc)
	if err != nil {
		return err
	}

	success := fmt.Sprintln("Build was successful")
	failure := fmt.Sprintln("Build failed")

	if cmdVal == 0 {
		note.WriteNote(
			tree,
			success,
		)
	} else {
		note.WriteNote(
			tree,
			failure,
		)
	}

	return nil
}
