package git_tmpdir

import (
	git "github.com/libgit2/git2go"
	"io/ioutil"
	"os"
	"time"
)

var sig = &git.Signature{
	Name:  "Test",
	Email: "test@asdf.com",
	When:  time.Now(),
}

func GitTmpDir(path string, cb func(*git.Repository)) error {
	path, err := ioutil.TempDir("", "tbd-capture-test")
	die(err)
	defer os.RemoveAll(path) // Cleanup on panic
	repo, err := git.InitRepository(path, false)
	die(err)
	bld, err := repo.TreeBuilder()
	die(err)
	blobid, err := repo.CreateBlobFromBuffer([]byte("This is a sample commit"))
	die(err)

	err = bld.Insert("README.MD", blobid, int(git.FilemodeBlob))
	die(err)
	treeId, err := bld.Write()
	die(err)
	tree, err := repo.LookupTree(treeId)
	die(err)
	commitId, err := repo.CreateCommit("", sig, sig, "Initial import", tree)
	die(err)
	// Force-create ref (override if not exists)
	_, err = repo.CreateReference("refs/heads/master", commitId, true, sig, "Initial commit")
	die(err)

	cb(repo)
	return os.RemoveAll(path) // Try to cleanup, and return any errors encountered.
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}
