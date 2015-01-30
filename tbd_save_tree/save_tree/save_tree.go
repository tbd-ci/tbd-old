package save_tree

import (
	"C"

	git "github.com/libgit2/git2go"
)

func Worktree() (string, error) {

	repo, err := git.OpenRepository(".")
	if err != nil {
		return "", err
	}

	tree, err := repo.RevparseSingle("HEAD^{tree}")
	if err != nil {
		return "", err
	}

	treeId := tree.Id()

	entry := git.IndexEntry{
		Mode: git.FilemodeTree,
		Id:   treeId,
		Path: ".",
	}

	index, err := git.NewIndex()
	if err != nil {
		return "", err
	}

	err = index.Add(&entry)
	if err != nil {
		return "", err
	}

	err = index.AddAll([]string{"."}, git.IndexAddDefault, nil)
	if err != nil {
		return "", err
	}

	treeOid, err := index.WriteTreeTo(repo)
	if err != nil {
		return "", err
	}

	return treeOid.String(), nil
}
