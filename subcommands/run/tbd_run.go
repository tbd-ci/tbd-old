package main

// TO DONT: Accept a worktree/ref (not the responsibility of this command)
// LATER: Sign build artefacts with your public key
// LATER: Windows

import (
	"flag"
	"fmt"
	git "github.com/libgit2/git2go"
	"github.com/tbd-ci/tbd/git/capture_output"
	"github.com/tbd-ci/tbd/git/empty_ref"
	"github.com/tbd-ci/tbd/git/nested_write"
	"os"
	"os/exec"
	"syscall"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("The final argument(s) are exec'd and the results stored in git.")
		fmt.Println("The working directory is expected to be the top level of a git repo")
	}

	config.refName = flag.String(
		"ref-name",
		"refs/tbd-ci-all-build-results",
		"(required): Store run output under this git ref",
	)
	config.storePrefix = flag.String(
		"store-prefix",
		"build",
		"Path to store build output under within the given ref",
	)
	config.repoDir = flag.String(
		"repo-dir",
		".git",
		"(required): Path to git repository to store results in",
	)
	config.propagateErrors = flag.Bool(
		"propagateErrors",
		true,
		"Exit with the same error code as the subcommand",
	)
}

type Config struct {
	refName         *string
	propagateErrors *bool
	repoDir         *string
	storePrefix     *string
}

func (c Config) RepoDir() string {
	if c.repoDir == nil {
		return ""
	} else {
		return *c.repoDir
	}
}

func (c Config) StorePrefix() string {
	if c.storePrefix == nil {
		return ""
	} else {
		return *c.storePrefix
	}
}

func (c Config) Ref() string {
	if c.refName == nil {
		return ""
	} else {
		return *c.refName
	}
}

func (c Config) Valid() bool {
	return !(c.Ref() == "" || c.RepoDir() == "" || (len(flag.Args()) == 0))
}

var config Config

func main() {
	flag.Parse()

	if !config.Valid() {
		flag.Usage()
		os.Exit(1)
	}

	repo, err := git.OpenRepository(config.RepoDir())
	if err != nil {
		fmt.Println("Error: Could not open git repository at ", config.RepoDir())
		os.Exit(1)
	}

	cmd := capture_output.Capture{
		Cmd:        exec.Command(flag.Args()[0], flag.Args()[1:]...),
		Repository: repo,
	}
	if cmd.Err() != nil {
		panic(err)
	}
	oid := cmd.Worktree()

	_, err = empty_ref.AssertRefIsCommit(repo, config.Ref(), nil)
	if err != nil {
		panic(err)
	}
	err = nested_write.AppendRef(
		nested_write.Paths{config.StorePrefix(): *oid},
		config.Ref(),
		repo,
	)
	if err != nil {
		panic(err)
	}

	// TODO: panics on systems without exit statuses (windows)
	waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)

	if *config.propagateErrors {
		os.Exit(waitStatus.ExitStatus())
	}
}
