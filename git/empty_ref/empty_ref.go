package empty_ref

import (
	git "github.com/libgit2/git2go"
	"time"
)

type CommitOptions struct {
	Author          *git.Signature
	Committer       *git.Signature
	ReflogWho       *git.Signature
	Message         string
	ReflogMessage   string
	InitialContents func(*git.TreeBuilder) error
}

func DefaultCommitOptions() *CommitOptions {
	author := git.Signature{
		Name:  "System",
		Email: "system@example.org",
		When:  time.Now(),
	}
	committer := author
	reflogWho := author
	return &CommitOptions{
		Author:        &author,
		Committer:     &committer,
		ReflogWho:     &reflogWho,
		Message:       "Initial Commit",
		ReflogMessage: "Generating ref",
	}
}

func AssertRefIsCommit(
	repo *git.Repository,
	refname string,
	opts *CommitOptions,
) (*git.Reference, error) {
	if ref, err := repo.LookupReference(refname); err == nil {
		target, err := repo.Lookup(ref.Target())
		if err == nil {
			if target.Type() == git.ObjectCommit {
				// It's already a commit; return it.
				return ref, nil
			}
		}
	}
	if opts == nil {
		opts = DefaultCommitOptions()
	}

	bld, err := repo.TreeBuilder()
	if err != nil {
		return nil, err
	}
	if opts.InitialContents != nil {
		err := opts.InitialContents(bld)
		if err != nil {
			return nil, err
		}
	}
	treeId, err := bld.Write()
	if err != nil {
		return nil, err
	}
	tree, err := repo.LookupTree(treeId)
	if err != nil {
		return nil, err
	}
	commitId, err := repo.CreateCommit(
		"", // Not writing to any ref yet
		opts.Author,
		opts.Committer,
		opts.Message,
		tree,
		// No parent
	)
	if err != nil {
		return nil, err
	}
	// Force-create ref (override if not exists)
	return repo.CreateReference(refname, commitId, true, opts.ReflogWho, opts.ReflogMessage)
}
