package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
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
		"",
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

func (c Config) Valid() bool {
	return !((c.prefix == nil) || (*c.prefix == "") || (len(flag.Args()) == 0))
}

var config Config

func main() {
	flag.Parse()

	if !config.Valid() {
		flag.Usage()
		os.Exit(1)
	}

	treeBeforeBuild := worktree()
	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	combined := bytes.Buffer{}

	cmd.Stdout = &multiBufferWriter{os.Stdout, &stdout, &combined}
	cmd.Stderr = &multiBufferWriter{os.Stderr, &stderr, &combined}

	run(cmd)

	// Will panic on systems without exit statuses
	waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)

	if *config.propagateErrors {
		defer func() {
			os.Exit(waitStatus.ExitStatus())
		}()
	}

	prefix := "refs/" + *config.prefix + "/" + treeBeforeBuild

	updateRef(prefix+"/STDOUT", hashFor(&stdout))
	updateRef(prefix+"/STDERR", hashFor(&stderr))
	updateRef(prefix+"/OUTPUT", hashFor(&combined))
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

func worktree() string {
	env := os.Environ()
	dir, err := ioutil.TempDir("", "tbd_index")
	die(err)
	env = append(env, "GIT_INDEX_FILE="+dir+"/index")

	output, err := withEnv(env, "git", "add", "-A", ".").CombinedOutput()

	die(err)

	worktree, err := withEnv(env, "git", "write-tree").Output()
	debug(worktree)
	die(err)

	return strip(string(worktree))
}

func withEnv(env []string, target ...string) *exec.Cmd {
	c := exec.Command(target[0], target[1:]...)
	c.Env = env
	return c
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
