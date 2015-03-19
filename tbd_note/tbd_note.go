package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/tbd-ci/tbd/tbd_note/note"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println("Creates build notes for a tree\n")
		fmt.Printf("%s [options...] <note message>\n", os.Args[0])
		flag.PrintDefaults()
	}

	// TODO: this default value won't work. Provide a sensible default
	tree := flag.String("tree", "HEAD^{tree}", "tree that note is attached to")
	config.tree = *tree
}

type Config struct {
	tree string
}

var config Config

func (c Config) valid() bool {
	// TODO: Use git to check for valid tree
	matched, err := regexp.MatchString(`[a-fA-F0-9]{40}`, config.tree)
	if err != nil {
		log.Fatal(err)
	}

	return matched
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 && config.valid() {
		flag.Usage()
		os.Exit(1)
	}

	noteId, err := note.WriteNote(config.tree, flag.Args()[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(noteId)
}
