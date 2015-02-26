package output

import (
	"log"

	git "github.com/libgit2/git2go"
)

func Display(treeish string) error {
	repo, err := git.OpenRepository(".")
	if err != nil {
		log.Fatal(err)
	}

	//treeOid, err := git.NewOid(treeish)
	//if err != nil {
	//  log.Fatal(err)
	//}

	object, err := repo.RevparseSingle(treeish)
	if err != nil {
		log.Fatal(err)
	}

	var treeOid *git.Oid

	switch object.Type() {
	case git.ObjectCommit:
		commit, err := repo.LookupCommit(object.Id())
		if err != nil {
			return err
		}

		tree, err := commit.Tree()
		if err != nil {
			return err
		}

		treeOid = tree.Id()

	case git.ObjectTree:
		treeOid = object.Id()
	default:
		log.Fatalf("%s is not a tag, commit or a tree object", treeish)
	}

	log.Println(treeOid.String())

	return nil
}
