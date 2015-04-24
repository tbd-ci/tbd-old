package capture_output

// TODO: Capture artifacts too (no, not in this package)

import (
	git "github.com/libgit2/git2go"
	"github.com/tbd-ci/tbd/git/tmpdir"
	"os/exec"
	"strings"
	"testing"
)

func testHarness(
	cmd *exec.Cmd,
	cb func(*git.Repository, *git.Tree),
) {
	err := git_tmpdir.GitTmpDir("tbd-capture-test", func(repo *git.Repository) {
		cap := Capture{
			Cmd:        cmd,
			Repository: repo,
		}
		treeOid := cap.Worktree()

		tree, err := repo.LookupTree(treeOid)
		if err != nil {
			panic(err)
		}
		cb(repo, tree)
	})
	if err != nil {
		panic(err)
	}
}

// Test that a process can have output captured
func TestCapture(t *testing.T) {
	runner := `
		echo 'stdout'
		echo 'stderr' 1>&2
		sleep 0.1
		echo 'stdout'
		echo 'stderr' 1>&2
	`
	testHarness(
		exec.Command("bash", "-c", runner),
		func(repo *git.Repository, tree *git.Tree) {
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
			contents := string(lookup(repo, tree.EntryByName("OUTPUT").Id).Contents())
			outs := strings.Count(contents, "stdout")
			if outs != 2 {
				t.Errorf(
					"Expected %s to have 2 copies of %s, had %d",
					contents,
					"stdout",
					outs,
				)
			}
			errs := strings.Count(contents, "stderr")
			if errs != 2 {
				t.Errorf(
					"Expected %s to have 2 copies of %s, had %d",
					contents,
					"stderr",
					errs,
				)
			}
		},
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
