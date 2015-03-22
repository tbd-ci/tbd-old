package output

import (
	"io/ioutil"
	"log"
	"path/filepath"

	git "github.com/libgit2/git2go"
)

type Output struct {
	Repo    *git.Repository
	Treeish string
	Build   string
	Stage   string
	Stream  string
}

func (c Output) filepath() string {
	return filepath.Join(
		String(c.treeFromTreeish()),
		c.Build,
		c.Stage,
		c.Stream,
	)
}

func (c Output) treeFromTreeish() (*git.Oid, error) {
	object, err := c.Repo.RevparseSingle(treeish)
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

func (c Output) readRef() error {

	return nil
}

func readBuild(path string) (buf []byte, err error) {
	buf, err = ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return
}
