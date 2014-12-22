package main

import (
	"flag"
	"fmt"
	git "github.com/libgit2/git2go"
	"os"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("The final argument(s) are exec'd and the results stored.")
		fmt.Println("The working directory is expected to be the top level of a git repo")
	}
	config.prefix = flag.String(
		"prefix",
		"refs/builds/",
		"(required): Store run output under this git ref prefix",
	)
	config.revision = flag.String(
		"revision",
		"master",
		"(required): Which commit or worktree to check the status of",
	)
	config.debug = flag.Bool(
		"debug",
		false,
		"Print verbose debugging",
	)
}

func debug(s interface{}) {
	if *config.debug {
		fmt.Println(s)
	}
}

type Config struct {
	prefix   *string
	revision *string
	debug    *bool
	repo     *git.Repository
}

func nilStr(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}

func (c Config) Prefix() string {
	return nilStr(c.prefix)
}

func (c Config) Target() git.Object {
	obj, err := c.Repo().RevparseSingle(nilStr(c.revision))
	die(err)
	return obj
}

func (c Config) Repo() *git.Repository {
	var err error
	if c.repo == nil {
		c.repo, err = git.OpenRepository(".")
		die(err)
	}
	return c.repo
}

func (c Config) Valid() bool {
	return (c.Prefix() != "")
}

var config Config

func main() {
	var err error
	flag.Parse()

	if !config.Valid() {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println(config.Target())
	rev := config.Target()
	// If a revision is given, use its tree
	if rev.Type() == git.ObjectCommit {
		rev, err = rev.(*git.Commit).Tree()
		die(err)
	}
	prefix := "refs/" + config.Prefix() + "/" + rev.Id().String() + "/STDOUT"
	ref, err := config.Repo().LookupReference(prefix)
	die(err)
	fmt.Println(ref)
}

func die(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, err.Error()+"\n")
	panic(err)
	os.Exit(1)
}
