package output

import "testing"

func TestReadFile(t *testing.T) {

	repo := git.OpenNewRepository(".")

	config := &Config{
		Repo:    repo,
		Treeish: "treeish",
		Stream:  "stream",
		Stage:   "stage",
		Build:   "build",
	}

	if config.filepath() != "treeish/build/stage/stream" {
		t.Fatalf("%s should have been %s", config.filepath, "treeish/build/stage/stream")
	}

}
