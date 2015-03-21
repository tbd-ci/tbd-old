package output

import (
	"io/ioutil"
	"log"
	"path/filepath"

	git "github.com/libgit2/git2go"
)

type Output struct {
	Repository *git.Repository
	Treeish    string
	Build      string
	Stage      string
	Stream     string
}

func (c Output) repo() error {
	if c.Repo == nil {
		repo, err := git.OpenRepository(".")
		if err != nil {
			return err
		}
		c.Repo = repo
	}
	return nil
}

func (c Output) filepath() string {
	return filepath.Join(c.Treeish, c.Build, c.Stage, c.Stream)
}

func (c Output) treeishOid() (*git.Oid, error) {
	treeOid, err := c.treeFromTreeish()
	if err != nil {
		return nil, err
	}

	return treeOid, nil
}

func (c Output) treeFromTreeish() (*git.Oid, error) {
	object, err := c.repo.RevparseSingle(treeish)
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

//func Display(string) error {
//  // TODO: handle -only flag
//  // TODO: handle "latest" for -build flag

//  repo, err := git.OpenRepository(".")
//  if err != nil {
//    log.Fatal(err)
//  }

//  treeOid, err := treeId(repo, treeish)
//  if err != nil {
//    log.Fatal(err)
//  }

//  treeId := treeOid.String()

//  outputPath := filepath.Join(
//    "refs",
//    "builds",
//    treeId,
//    config.Build,
//    config.Stream,
//  )

//  buf, err := readBuild(outputPath)
//  if err != nil {
//    log.Fatal(err)
//  }

//  blob, err := buf.Read()
//  if err != nil {
//    log.Fatal(err)
//  }

//  fmt.Println(blob)

//  return nil
//}

func readBuild(path string) (buf []byte, err error) {
	buf, err = ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return
}
