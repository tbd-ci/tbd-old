package mktmp

import (
	"io/ioutil"

	"github.com/libgit2/git2go"
)

func CheckoutTmp(tree string) (string, error) {
	tmpPath, err := ioutil.TempDir("", "tbd-build")
	if err != nil {
		return "", err
	}

	repo, err := git.OpenRepository(".")
	if err != nil {
		return "", err
	}

	treeOid, err := git.NewOid(tree)
	if err != nil {
		return "", err
	}

	treeId, err := repo.LookupTree(treeOid)
	if err != nil {
		return "", err
	}

	var checkoutOpts git.CheckoutOpts

	checkoutOpts.Strategy = git.CheckoutForce
	checkoutOpts.TargetDirectory = tmpPath

	err = repo.CheckoutTree(treeId, &checkoutOpts)
	if err != nil {
		return "", err
	}

	return tmpPath, nil
}
