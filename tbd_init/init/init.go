package init

import (
	"log"
	"os"
)

func Init() error {
	if _, err := os.Open("./ci"); err != nil {
		log.Fatal(err)
	}
	return nil
}
