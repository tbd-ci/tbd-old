package main

import (
	"regexp"
	"testing"

	"github.com/libgit2/git2go"
)

func TestTmpDir(t *testing.T) {
	tmpdir, err := mktmp_d()
	if err != nil {
		t.Fatal("error with mktmp_d")
	}

	match, err := regexp.MatchString("/tmp/tmp.[a-zA-Z0-9]{10}", tmpdir)
	if err != nil {
		t.Fatal("Error with matching")
	}

	if !match {
		t.Fatalf("%s was not even close to /tmp/tmp.2Kjy3E8hw8", tmpdir)
	}
}

//func TestCheckoutTmp(t *testing.T) {
//  repo := git.createTestRepo(t)
//  defer os.RemoveAll(repo.Workdir())

//  idx, err := repo.Index()
//  checkFatal(t, err)
//  err = idx.AddByPath("README")
//  checkFatal(t, err)
//  treeId, err := idx.WriteTree()
//  checkFatal(t, err)

//  treeString := string(treeId)

//  path, err := checkout_tmp(&treeString)
//  checkFatal(t, err)

//}

//func checkFatal(t *testing.T, err error) {
//  if err == nil {
//    return
//  }

//  // The failure happens at wherever we were called, not here
//  _, file, line, ok := runtime.Caller(1)
//  if !ok {
//    t.Fatal()
//  }

//  t.Fatalf("Fail at %v:%v; %v", file, line, err)
//}
