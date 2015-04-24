package save_tree

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func strip(s string) string {
	return strings.Replace(string(s), "\n", "", -1)
}

func Worktree() (string, error) {
	// TODO: Handle the various errors which git can produce.
	// Currently this just prints 'Error: Exit code <x>'.
	// Probably want to show the error from git (e.g. no `.git` dir found)
	env := os.Environ()
	dir, err := ioutil.TempDir("", "tbd_index")
	if err != nil {
		return "", err
	}
	env = append(env, "GIT_INDEX_FILE="+dir+"/index")

	_, err = withEnv(env, "git", "add", "-A", ".").CombinedOutput()

	if err != nil {
		return "", err
	}

	worktree, err := withEnv(env, "git", "write-tree").Output()
	if err != nil {
		return "", err
	}

	return strip(string(worktree)), nil
}

func withEnv(env []string, target ...string) *exec.Cmd {
	c := exec.Command(target[0], target[1:]...)
	c.Env = env
	return c
}
