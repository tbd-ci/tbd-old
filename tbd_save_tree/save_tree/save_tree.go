package save_tree

import (
	"C"
	"os"

	git "github.com/libgit2/git2go"
)

func Worktree() (string, error) {

	workdir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	repo, err := git.OpenRepository(".")
	if err != nil {
		return "", err
	}

	index, err := repo.Index()
	if err != nil {
		return "", err
	}

	err = index.AddAll([]string{workdir}, git.IndexAddDefault, nil)
	if err != nil {
		return "", err
	}

	treeOid, err := index.WriteTree()
	if err != nil {
		return "", err
	}

	return treeOid.String(), nil
}
