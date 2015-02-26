package output

import (
	"io/ioutil"
	"log"
	"path/filepath"

	git "github.com/libgit2/git2go"
)

func Display(config *Config, treeish string) error {
	// TODO: handle -only flag
	// TODO: handle "latest" for -build flag
	repo, err := git.OpenRepository(".")
	if err != nil {
		log.Fatal(err)
	}

	treeOid, err := treeId(repo, treeish)
	if err != nil {
		log.Fatal(err)
	}

	treeId := treeOid.String()

	return nil
}

func readBuild(ref, treeOid string) ([]byte, error) {
	buf, err := ioutil.ReadFile(filepath.Join("refs", "builds", treeOid, stage, stream))
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}

func treeId(repo *git.Repository, treeish string) (*git.Oid, error) {
	object, err := repo.RevparseSingle(treeish)
	if err != nil {
		log.Fatal(err)
	}

	var treeOid *git.Oid

	switch object.Type() {
	case git.ObjectCommit:
		commit, err := repo.LookupCommit(object.Id())
		if err != nil {
			return nil, err
		}

		tree, err := commit.Tree()
		if err != nil {
			return nil, err
		}

		treeOid = tree.Id()

	case git.ObjectTree:
		treeOid = object.Id()
	default:
		log.Fatalf("%s is not a tag, commit or a tree object", treeish)
	}

	return treeOid, nil
}
