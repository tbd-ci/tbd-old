package main

// TODO: should this use git2go instead of CLI?
// TODO: if not using git2go, extract types for each git subcommand
// TODO: Save exit status
// TODO: Save worktree
// TODO: Save timestamped versions of the build output
// TODO: Make prefix default to blank. Caller should put the worktree into the prefix.
// TO DONT: Accept a worktree/ref (not the responsibility of this task)
// LATER: Sign build artefacts with your public key
// LATER: Windows

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("The final argument(s) are exec'd and the results stored in git.")
		fmt.Println("The working directory is expected to be the top level of a git repo")
	}

	config.prefix = flag.String(
		"prefix",
		"refs/builds/",
		"(required): Store run output under this git ref prefix",
	)
	config.propagateErrors = flag.Bool(
		"propagateErrors",
		true,
		"Exit with the same error code as the subcommand",
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
	prefix          *string
	propagateErrors *bool
	debug           *bool
}

func (c Config) Prefix() string {
	if c.prefix == nil {
		return ""
	} else {
		return *c.prefix
	}
}

func (c Config) Valid() bool {
	return !(c.Prefix() == "" || (len(flag.Args()) == 0))
}

var config Config

func main() {
	flag.Parse()

	if !config.Valid() {
		flag.Usage()
		os.Exit(1)
	}

	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	combined := bytes.Buffer{}

	cmd.Stdout = &multiBufferWriter{os.Stdout, &stdout, &combined}
	cmd.Stderr = &multiBufferWriter{os.Stderr, &stderr, &combined}

	run(cmd)

	// TODO: panics on systems without exit statuses (windows)
	waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)

	if *config.propagateErrors {
		defer func() {
			// TODO: This will exit with success if the build passed
			// but updating the ref paniced
			// die() is probably a bad idea.
			os.Exit(waitStatus.ExitStatus())
		}()
	}

	updateRef(*config.prefix+"/STDOUT", hashFor(&stdout))
	updateRef(*config.prefix+"/STDERR", hashFor(&stderr))
	updateRef(*config.prefix+"/OUTPUT", hashFor(&combined))
}

func run(cmd *exec.Cmd) {
	die(cmd.Start())
	err := cmd.Wait()
	if err != nil {
		// ExitError indicates process had a non-zero exit code.
		// This is expected behavior for us.
		_, isExitErr := err.(*exec.ExitError)
		if !isExitErr {
			die(err)
		}
	}
}

func hashFor(r io.Reader) string {
	cmd := exec.Command("git", "hash-object", "-w", "--stdin")
	cmd.Stdin = r
	out, err := cmd.Output()
	die(err)
	return strip(string(out))
}

func strip(s string) string {
	return strings.Replace(string(s), "\n", "", -1)
}

func updateRef(ref, sha string) {
	cmd := exec.Command("git", "update-ref", ref, sha)
	out, err := cmd.CombinedOutput()
	debug(string(out))
	die(err)
}

type multiBufferWriter struct {
	primary  *os.File
	output   *bytes.Buffer
	combined *bytes.Buffer
}

func (db *multiBufferWriter) Write(p []byte) (int, error) {
	n, err := db.primary.Write(p)
	// bytes.Buffer never returns an error and always writes the full length.
	// We can safely discard its result.
	db.output.Write(p[:n])
	db.combined.Write(p[:n])
	return n, err
}

func die(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, err.Error()+"\n")
	panic(err)
	os.Exit(1)
}
