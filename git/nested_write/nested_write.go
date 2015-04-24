package nested_write

import (
	git "github.com/libgit2/git2go"
	"github.com/tbd-ci/tbd/git/empty_ref"
	"time"
)

type Paths map[string]git.Oid

func AppendTreeViaIndex(p Paths, t *git.Tree) (*git.Tree, error) {
	idx, err := t.Owner().Index()
	if err != nil {
		return nil, err
	}
	for path, oid := range p {
		err = idx.Add(&git.IndexEntry{
			Ctime: time.Now(),
			Mtime: time.Now(),
			Mode:  git.FilemodeBlob,
			Uid:   0,
			Gid:   0,
			Size:  0,
			Id:    &oid,
			Path:  path,
		})
		if err != nil {
			return nil, err
		}
	}
	treeOid, err := idx.WriteTree()
	if err != nil {
		return nil, err
	}
	return t.Owner().LookupTree(treeOid)
}

func AppendRef(p Paths, ref string, in *git.Repository, author *git.Signature) error {
	commit, err := Append(p, ref, in, author)
	if err != nil {
		return err
	}
	refb, err := in.LookupReference(ref)
	if err != nil {
		// This ref doesn't exist; create it.
		commitOpts := empty_ref.DefaultCommitOptions()
		commitOpts.Author = commit.Author()
		commitOpts.Committer = commit.Committer()
		commitOpts.ReflogWho = commit.Committer()

		refb, err = empty_ref.AssertRefIsCommit(
			in,
			ref,
			commitOpts,
		)
		if err != nil {
			return err
		}
	}
	refb.SetTarget(commit.Id(), commit.Author(), "Completed build")
	return nil
}

func Append(p Paths, ref string, in *git.Repository, author *git.Signature) (*git.Commit, error) {
	committish := Lookup(in, ref)
	commits := []*git.Commit{committish.commit}
	if committish.err != nil {
		// ref doesn't exist; use a nil commit and an empty tree.
		commits = []*git.Commit{}
		bld, err := in.TreeBuilder()
		if err != nil {
			return nil, err
		}
		treeId, err := bld.Write()
		if err != nil {
			return nil, err
		}
		committish.tree, err = in.LookupTree(treeId)
		if err != nil {
			return nil, err
		}
	}

	tree, err := AppendTreeViaIndex(p, committish.tree)
	if err != nil {
		return nil, err
	}
	commitId, err := in.CreateCommit(
		"",
		author,
		author,
		"Add build",
		tree,
		commits...,
	)
	if err != nil {
		return nil, err
	}
	return in.LookupCommit(commitId)
}

type Committish struct {
	err    error
	ref    *git.Reference
	commit *git.Commit
	tree   *git.Tree
}

func Lookup(repo *git.Repository, ref string) (c Committish) {
	c.ref, c.err = repo.LookupReference(ref)
	if c.err != nil {
		return
	}
	c.commit, c.err = repo.LookupCommit(c.ref.Target())
	if c.err != nil {
		return
	}
	c.tree, c.err = c.commit.Tree()
	if c.err != nil {
		return
	}
	return
}
