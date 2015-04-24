package note

import (
	"fmt"
	"time"

	git "github.com/libgit2/git2go"
)

type User struct {
	name  string
	email string
}

func userDetails() (User, error) {
	var user User

	configPath, err := git.ConfigFindGlobal()
	if err != nil {
		return user, err
	}

	config, err := git.OpenOndisk(nil, configPath)
	if err != nil {
		return user, err
	}

	user.name, err = config.LookupString("user.name")
	if err != nil {
		return user, err
	}
	user.email, err = config.LookupString("user.email")
	if err != nil {
		return user, err
	}

	return user, err
}

func WriteNote(tree, buildResult string) (string, error) {
	repo, err := git.OpenRepository(".")
	if err != nil {
		return "", err
	}

	treeOid, err := git.NewOid(tree)
	if err != nil {
		return "", err
	}

	var authorSig git.Signature

	user, err := userDetails()
	if err != nil {
		return "", err
	}

	authorSig.Name = user.name
	authorSig.Email = user.email
	authorSig.When = time.Now()

	// ref string, author, committer *Signature, id *Oid,
	// note string, force bool
	noteId, err := repo.CreateNote(
		"",
		&authorSig,
		&authorSig,
		treeOid,
		fmt.Sprintf("%s at %v", buildResult, time.Now()),
		true,
	)
	if err != nil {
		return "", err
	}

	return noteId.String(), nil
}
