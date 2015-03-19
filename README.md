# tbd

A decentralised ci built on git.

## Status

Pre-Alpha. You probably don't want to use this yet.

## Why a new CI?

We're frustrated with

 * The exorbitant cost of paid solutions (*cough* bamboo)
 * Inflexible CI servers making it difficult to
  - customise your workflow (e.g. merge master into branch, then build - as the `rust` project does)
  - port build configuration between systems (or even a different server running the same software)
  - answer simple questions like `did commit <sha> 2 months ago pass or fail?`
  - answer simple questions like `What was the build configuration when <commit> failed the build?`
 * Workflows which encourage developers not to check the build status of their commits
 * Tools that *only* provide web interfaces
 * CIs claiming flexibility because you can upload a plugin
  - Plugins are separately versioned & frequently implemented in a different language
  - It's not clear which version of a plugin was running when a build ran

## How is tbd different?

Core ideas driving `tbd` which are different to traditional CI:

 * All CI configuration is maintained in the same repository as your code
 * Build results are stored in your git repository (not in the main tree)as git notes, alongside that commit's code

## Build storage

`tbd` stores build results in a git branch (called `tbd-ci` by default, controlled by TBD_STORAGE_BRANCH).
When a build is triggered by a commit, the results are written to `tbd-ci` as a merge commit. If that tree has already been built, you'll still get a merge commit but it will have no changes.
If a build is triggered by a worktree you just get a regular commit; there's no second parent for the merge.

`tbd-ci` commits contain a directory for each commit/worktree which has ever been built.
Because of how git stores files, this requires very little storage.

Example directory structure:

```
<source worktree sha>
  <spec:coverage>
    <build timestamp and host>
      STDOUT
      STDERR
      ETC
      <tree after build>
        application files
        build artifacts
<source commit sha>
  <same thing>
```

When you run a build, you create a new commit and update the `tbd-ci` branch

To check a file that was modified by the build process:
`git show tbd-ci:<commit>/<metadata-hash>/<target>/WORKTREE/<artifact>`
`tbd show <ref-like> <target> [--build-number <metadata-hash, defaults to most recent>] WORKTREE/<artifact>`

### Advantages
 * Builds can be viewed & data extracted without tbd tools installed
 * no changes required to repository config
 * git notes for each commit can link directly to the build

### Disadvantages
No garbage collection for builds of a tree that was never pushed
 * not a problem unless devs are doing local builds and pushing the branch
 * probably not a problem since even on a large repo there's very little change
 * We could write a tool to strip out builds which aren't in the history (it'd re-write the branch)

It's easy commit a large binary as part of the post-build worktree and hard to undo.
 * We'll need to make it obvious, when writing a build to git, that artifacts have been saved (name/size)
A developer could carelessly check out `tbd-ci`, which would cause a *lot* of files to be written to their machine.
 * If you're concerned about this, configure tbd to write to a ref instead
 * TODO: Should we always write to a ref and require developers to configure their repo appropriately?
