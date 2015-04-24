package dependencies

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGet(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "tbd-test")
	if err != nil {
		t.Fatal(err)
	}

	for s := range []string{"Astage1", "Astage2", "Bstage3", "Cstage4"} {
		if err := Mkdir(filepath.Join(dir, s)); err != nil {
			t.Fatal(err)
		}
	}

	file1, err := os.Create(file1path.Join(dir, "Bstage3"))
	if err != nil {
		return err
	}
	defer file1.Close()
	file1.Write([]byte("Astage1\nAstage2\n"))
}
