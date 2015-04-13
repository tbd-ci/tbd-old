// Meta-package to specify global dependencies.
package tbd

import (
	_ "github.com/libgit2/git2go"
	_ "github.com/tbd-ci/tbd/capture_output"
)

const cmd = `

# git tbd show <source> <target> --build-number 3

exec(git show tbd-ci-all-builds -- <source>/<target>/<buildname>)

`

// beta test: atlas

// 1) create ref by hand
// 2) run spec:coverage and record stdout/stderr/combined into a tmpdir (a single command; tbd-write-build)
// 3) write that tree into the ref
// 4) explore it using hand-coded everything

// cd tmpdir
// checkout repo
// GIT_DIR='~/atlas/.git' tbd-write-build ci/spec:coverage/build
//   -> stdout | git put-object
//   -> stderr | git put-object
//   -> git write-tree <custom tree>
//   -> output sha to the fresh tree
// update ref `tbd-ci-all-build-results`:
// <sourceRef> -> <spec:coverage> -> <host> -> <output of tbd-write-build>
