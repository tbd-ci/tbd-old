package init

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Init() error {
	if _, err := os.Open("./ci"); err != nil {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("You do not have a ci/ directory. Create one now? (y/n)")

		text, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if strings.Trim(text, "\n\r") == "y" || strings.Trim(text, "\n\r") == "Y" {
			if err := os.Mkdir("ci", os.ModeDir); err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println(text)
	}
	return nil
}
