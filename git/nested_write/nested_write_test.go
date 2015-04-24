package nested_write

// TODO: Test writing a tree as well as a blob.
// TODO: Test writing to a subtree which already has content

import (
	git "github.com/libgit2/git2go"
	"github.com/tbd-ci/tbd/git/tmpdir"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var sig = git.Signature{
	Name:  "Test",
	Email: "test@asdf.com",
	When:  time.Now(),
}

func TestNestedWrite(t *testing.T) {
	git_tmpdir.GitTmpDir("nested-write", func(repo *git.Repository) {
		file, err := repo.CreateBlobFromBuffer([]byte("file contents"))
		die(err)
		paths := Paths(make(map[string]git.Oid, 2))
		paths["first/second/third.go"] = *file

		author := git.Signature{
			Name:  "TBD",
			Email: "tbd@example.org",
			When:  time.Unix(0, 0),
		}

		err = AppendRef(paths, "refs/foobar", repo, &author)
		die(err)
		committish := Lookup(repo, "refs/foobar")
		tree, err := committish.tree, committish.err
		die(err)
		if tree.EntryCount() < 1 {
			t.Error("Expected tree to have entries")
			return
		}
		if tree.EntryByIndex(0).Name != "first" {
			t.Errorf("Expected first entry to be named 'first', got %s", tree.EntryByIndex(0).Name)
			return
		}
		entry, err := tree.EntryByPath("first")
		die(err)
		if entry.Type != git.ObjectTree {
			t.Error("Expected first to be a folder")
			return
		}
		if entry.Id.String() != "46736b97f44a7d10e92dadfdfc5c7a921f191803" {
			t.Errorf("Expected first id to be %s", entry.Id.String())
			return
		}
		subtree, err := repo.LookupTree(entry.Id)
		die(err)

		entry, err = subtree.EntryByPath("second")
		if entry.Type != git.ObjectTree {
			t.Error("Expected second to be a folder")
			return
		}
		if entry.Id.String() != "87e03b41197ce32e98342a393da887683544a712" {
			t.Errorf("Expected second id to be %s", entry.Id.String())
			return
		}
		subtree, err = repo.LookupTree(entry.Id)
		die(err)
		entry, err = subtree.EntryByPath("third.go")
		die(err)
		if entry.Type != git.ObjectBlob {
			t.Error("Expected third to be a blob")
			return
		}
		if entry.Id.String() != "754bb844fb01df2613c0c1fe26eaa701ce46e853" {
			t.Errorf("Expected third id to be %s", entry.Id.String())
			return
		}
		blob, err := repo.LookupBlob(entry.Id)
		die(err)
		if string(blob.Contents()) != "file contents" {
			t.Errorf("Expected blob contents to be 'file contents', got '%s'", string(blob.Contents()))
		}

		revParseOut, err := runGitPath(repo.Path(), "git", "rev-parse", "refs/foobar")
		dieOf(err, revParseOut)
		if revParseOut != committish.commit.Id().String() {
			t.Errorf("Expected output of rev-parse to be 'c5d38cdeb4f6e98ee3646d0b73d65734a9ac4596', got '%s'", revParseOut)
		}

		rawTree, err := runGitPath(repo.Path(), "git", "cat-file", "-p", committish.tree.Id().String())
		dieOf(err, rawTree)
		expectedTree := `040000 tree 46736b97f44a7d10e92dadfdfc5c7a921f191803	first`
		if rawTree != expectedTree {
			t.Errorf("Expected output of cat-file to be '%s', got '%s'", expectedTree, rawTree)
		}
	})
}

func runGitPath(path string, command ...string) (output string, err error) {
	c := exec.Command(command[0], command[1:]...)
	c.Env = append(c.Env, "GIT_DIR="+path)
	b, err := c.CombinedOutput()
	output = strings.Trim(string(b), " \n")
	return
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func dieOf(err error, reason string) {
	if err != nil {
		panic(err.Error() + " " + reason)
	}
}
