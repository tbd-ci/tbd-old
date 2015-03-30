// Runs a command and captures its stdout and stderr
package capture_output

import (
	"bytes"
	git "github.com/libgit2/git2go"
	"io"
	"os"
	"os/exec"
)

type Capture struct {
	*exec.Cmd
	*git.Repository
}

func (c Capture) Worktree() (*git.Oid, error) {
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

	err := run(cmd)
	if err != nil {
		return nil, err
	}

	tree, err := c.Repository.TreeBuilder()
	if err != nil {
		return nil, err
	}
	err = c.writeTree(tree, "STDOUT", stdout.Bytes())
	if err != nil {
		return nil, err
	}
	err = c.writeTree(tree, "STDERR", stderr.Bytes())
	if err != nil {
		return nil, err
	}
	err = c.writeTree(tree, "OUTPUT", combined.Bytes())
	if err != nil {
		return nil, err
	}
	treeOid, err := tree.Write()
	if err != nil {
		return nil, err
	}

	return treeOid, nil
}

func (c Capture) writeTree(tree *git.TreeBuilder, path string, bytes []byte) error {
	blobid, err := c.Repository.CreateBlobFromBuffer(bytes)
	if err != nil {
		return err
	}
	return tree.Insert(path, blobid, int(git.FilemodeBlob))
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

func run(cmd *exec.Cmd) error {
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		// ExitError indicates process had a non-zero exit code.
		// This is expected behavior for us.
		_, isExitErr := err.(*exec.ExitError)
		if !isExitErr {
			return err
		}
	}
	return nil
}
