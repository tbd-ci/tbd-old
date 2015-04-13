// TODO: Tests
package lookup

import (
	git "github.com/libgit2/git2go"
)

type Treeish struct {
	tree   *git.Tree
	err error
}

func NewTreeishOid(repo *git.Repository, id *git.Oid) Treeish {
	obj, err := repo.Lookup(id)
	if err != nil {
		return Treeish{err: err}
	}
	return NewTreeishObject(repo, obj)
}

type targeter interface {
	Target() (*git.Oid, error)
}

func NewTreeishObject(repo *git.Repository, obj git.Object) Treeish {
	switch t := target.(type) {
	case *git.Blob:
		return Treeish{err: fmt.Errorf("%s is a blob, not a tree.", obj.Id())}
	case *git.Commit:
		target, err := repo.Lookup(t.TreeId())
		if err != nil {
			return Treeish{err: error}
		}
    return NewTreeishObject(repo, target)
	case *git.Tree:
		return Treeish{
			tree: t,
      repo: repo
		}
	case targeter: // Handles tag, branch, ref
		target, err := repo.Lookup(t.Target())
		if err != nil {
			return Treeish{err: error}
		}
		return NewTreeishObject(repo, target)
	default:
		return Treeish{err: fmt.Errorf("%s is not an object.", name)}
	}
}

func NewTreeish(repo *git.Repository, name string) Treeish {
	target, err := repo.RevparseSingle(name)
	if err != nil {
		return Treeish{err: err}
	}
	return NewTreeishObject(repo, target)
}

func (t Treeish) Tree() *git.Tree {
  return t.tree
}

func (t Treeish) Err() error {
  return t.err
}
