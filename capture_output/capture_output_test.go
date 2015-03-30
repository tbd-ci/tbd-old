package capture_output

// TODO: Capture artifacts too

import (
	git "github.com/libgit2/git2go"
	"io/ioutil"
	"os/exec"
	"testing"
)

func testHarness(t *testing.T, cmd *exec.Cmd) (*git.Repository, *git.Tree) {
	path, err := ioutil.TempDir("", "tbd-capture-test")
	if err != nil {
		panic(err)
	}
	_, err = exec.Command("git", "init", path).CombinedOutput()
	if err != nil {
		panic(err)
	}
	repo, err := git.InitRepository(path, false)
	if err != nil {
		panic(err)
	}

	// Test that a process can have output & captured
	cap := Capture{
		cmd,
		repo,
	}
	treeOid, err := cap.Worktree()
	if err != nil {
		panic(err)
	}

	tree, err := repo.LookupTree(treeOid)
	if err != nil {
		panic(err)
	}
	return repo, tree
}

func TestCapture(t *testing.T) {
	runner := `
		echo 'stdout'
		echo 'stderr' 1>&2
		sleep 0.1
		echo 'stdout'
		echo 'stderr' 1>&2
	`
	repo, tree := testHarness(t, exec.Command("bash", "-c", runner))

	eq(
		t,
		string(lookup(repo, tree.EntryByName("STDOUT").Id).Contents()),
		"stdout\nstdout\n",
	)
	eq(
		t,
		string(lookup(repo, tree.EntryByName("STDERR").Id).Contents()),
		"stderr\nstderr\n",
	)

	eq(
		t,
		string(lookup(repo, tree.EntryByName("OUTPUT").Id).Contents()),
		"stdout\nstderr\nstdout\nstderr\n",
	)
}

func lookup(repo *git.Repository, id *git.Oid) *git.Blob {
	combined, err := repo.LookupBlob(id)
	if err != nil {
		panic(err)
	}
	return combined
}
func eq(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("Expected '%s' to be '%s'", actual, expected)
	}
}
