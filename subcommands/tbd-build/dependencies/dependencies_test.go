package dependencies

import "testing"

func TestNextBatch(t *testing.T) {
	depTree := make(map[string][]string)

	// dependencies a -> b --- a depends on b
	// c -> a
	// c -> b
	// d -> a
	// d -> b
	// e -> c
	// f -> d
	// g -> e
	// h -> f
	// h -> e
	// i -> g
	// i -> h

	depTree["taskA"] = []string{}
	depTree["taskB"] = []string{}
	depTree["taskC"] = []string{"taskA", "taskB"}
	depTree["taskD"] = []string{"taskA", "taskB"}
	depTree["taskE"] = []string{"taskC"}
	depTree["taskF"] = []string{"taskD"}
	depTree["taskG"] = []string{"taskE"}
	depTree["taskH"] = []string{"taskE", "taskF"}
	depTree["taskI"] = []string{"taskG", "taskH"}

	deps := Dependencies{
		dependent: depTree,
	}

	next, err := NextBatch(&deps)
	if err != nil {
		t.Fatal(err)
	}

	if next != "taskA" {
		t.Errorf("%v should have been taskA", next)
	}

	next, err = NextBatch(&deps)
	if err != nil {
		t.Fatal(err)
	}

	// taskA complete, remove from depTree
	if next != "taskB" {
		t.Errorf("%v should have been [taskA, taskB]", next)
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
