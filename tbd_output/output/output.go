package output

import (
	"fmt"
	"log"
	"path/filepath"

	git "github.com/libgit2/git2go"
)

type Config struct {
	Build  string
	Stage  string
	Stream string
}

func Display(treeish string, config Config) error {
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

	outputPath := filepath.Join(
		"refs",
		"builds",
		treeId,
		config.Build,
		config.Stream,
	)

	buf, err := readBuild(outputPath)
	if err != nil {
		log.Fatal(err)
	}

	blob, err := buf.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(blob)

	return nil
}

func readBuild(path string) (buf []byte, err error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return
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
