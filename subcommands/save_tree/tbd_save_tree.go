package main

import (
	"fmt"
	"github.com/tbd-ci/tbd/subcommands/save_tree/save_tree"
	"os"
)

func main() {
	tree, err := save_tree.Worktree()
	if err != nil {
		// TODO: This thing returns *terrible* errors.
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println(tree)
}
