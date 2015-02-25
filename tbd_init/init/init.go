package init

import (
	"bufio"
	"fmt"
	"os"
)

func Init() error {
	if _, err := os.Open("./ci"); err != nil {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("You do not have a ci/ directory. Create one now? (y/n)")

		text, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		fmt.Println(text)
	}
	return nil
}
