// Runs a command and captures its stdout and stderr
// TODO: Save exit status
// TODO: Save timestamped versions of the build output
// TODO: Save worktree
package capture_output

import (
	"bytes"
	"fmt"
	git "github.com/libgit2/git2go"
	"io"
	"os"
	"os/exec"
)

type Capture struct {
	*exec.Cmd
	*git.Repository

	finished bool
	err      error
	stdout   *git.Oid
	stderr   *git.Oid
	combined *git.Oid
	tree     *git.Oid
}

func (c *Capture) run() {
	if c.finished {
		return
	}
	cmd := c.Cmd
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	combined := bytes.Buffer{}

	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}
	cmd.Stdout = &multiBufferWriter{cmd.Stdout, &stdout, &combined}
	cmd.Stderr = &multiBufferWriter{cmd.Stderr, &stderr, &combined}

	c.err = cmd.Start()
	if c.err != nil {
		return
	}
	err := cmd.Wait()
	if err != nil {
		// ExitError indicates process had a non-zero exit code.
		// This is expected behavior for us.

		_, isExitErr := err.(*exec.ExitError)
		if !isExitErr {
			c.err = err
			return
		}
	}

	c.stdout, c.err = c.Repository.CreateBlobFromBuffer(stdout.Bytes())
	if c.err != nil {
		return
	}
	c.stderr, c.err = c.Repository.CreateBlobFromBuffer(stderr.Bytes())
	if c.err != nil {
		return
	}
	c.combined, c.err = c.Repository.CreateBlobFromBuffer(combined.Bytes())
	if c.err != nil {
		return
	}

	var tree *git.TreeBuilder
	tree, c.err = c.Repository.TreeBuilder()
	if c.err != nil {
		return
	}
	c.err = tree.Insert("STDOUT", c.stdout, int(git.FilemodeBlob))
	if c.err != nil {
		return
	}
	c.err = tree.Insert("STDERR", c.stderr, int(git.FilemodeBlob))
	if c.err != nil {
		return
	}
	c.err = tree.Insert("OUTPUT", c.combined, int(git.FilemodeBlob))
	if c.err != nil {
		return
	}
	c.tree, c.err = tree.Write()
	if c.err != nil {
		return
	}
	fmt.Println(c.tree.String())
	c.finished = true
}

func (c *Capture) Err() error {
	c.run()
	return c.err
}

func (c *Capture) Worktree() *git.Oid {
	c.run()
	return c.tree
}

type multiBufferWriter struct {
	primary  io.Writer
	output   *bytes.Buffer
	combined *bytes.Buffer
}

func (db *multiBufferWriter) Write(p []byte) (int, error) {
	n, err := db.primary.Write(p)
	// bytes.Buffer never returns an error and always writes the full length.
	// We can safely discard its result.
	db.output.Write(p[:n])
	db.combined.Write(p[:n])
	return n, err
}
