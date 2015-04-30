package dependencies

import "testing"

func TestNext(t *testing.T) {
	depTree := make(map[string][]string)

	// dependencies a -> b --- a depends on b
	// c -> a
	// d -> b
	// e -> c
	// e -> d

	depTree["taskA"] = []string{}
	depTree["taskB"] = []string{}
	depTree["taskC"] = []string{"taskA"}
	depTree["taskD"] = []string{"taskB"}
	depTree["taskE"] = []string{"taskC", "taskD"}

	deps := Dependencies{
		dependent: depTree,
	}

	next, err := Next(&deps)
	if err != nil {
		t.Fatal(err)
	}

	if next != "taskA" && next != "taskB" {
		t.Errorf("%v should have been taskA or taskB", next)
	}
}

//func TestGet(t *testing.T) {
//  dir, err := ioutil.TempDir("/tmp", "tbd-test")
//  if err != nil {
//    t.Fatal(err)
//  }

//  for s := range []string{"Astage1", "Astage2", "Bstage3", "Cstage4"} {
//    if err := Mkdir(filepath.Join(dir, s)); err != nil {
//      t.Fatal(err)
//    }
//  }

//  file1, err := os.Create(file1path.Join(dir, "Bstage3"))
//  if err != nil {
//    return err
//  }
//  defer file1.Close()
//  file1.Write([]byte("Astage1\nAstage2\n"))
//}
