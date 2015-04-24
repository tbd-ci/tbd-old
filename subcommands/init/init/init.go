package init

import (
	"fmt"
	"log"
	"os"
)

func Init() error {
	exists, err := pathExists("ci")
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		fmt.Print("You already have a ci/ directory.",
			"\n",
			"If this is not for tbd you will need to configure tbd to use an alternate directory (see man tbd_build)\n",
		)
	} else {
		fmt.Println("Adding `./ci` directory...")
		if err := os.Mkdir("ci", 0775); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat("ci")
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	} else {
		log.Fatal(err)
	}
	return true, err
}
