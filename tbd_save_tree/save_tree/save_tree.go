// #include <index.h>
package save_tree

import (
	"C"

	git "github.com/libgit2/git2go"
)

var afterIndexCreated = func(string, string) int {
	return 0
}

func Worktree() (string, error) {
	// TODO: Handle the various errors which git can produce.
	// Currently this just prints 'Error: Exit code <x>'.
	// Probably want to show the error from git (e.g. no `.git` dir found)

	repo, err := git.OpenRepository(".")
	if err != nil {
		return "", err
	}

	index, err := git.NewIndex()
	if err != nil {
		return "", err
	}

	index.AddAll([]string{"."}, git.GIT_INDEX_ADD_DEFAULT, afterIndexCreated)

	treeOid, err := index.WriteTreeTo(repo)
	if err != nil {
		return "", err
	}

	return treeOid.String(), nil
}
