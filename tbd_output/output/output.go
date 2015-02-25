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

	treeOid, err := git.NewOid(treeish)
	if err != nil {
		log.Fatal(err)
	}

	object, err := repo.Lookup(treeOid)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(object.Type())

	return nil
}
