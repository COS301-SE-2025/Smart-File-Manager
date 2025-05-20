package main

import (
	"fmt"
	"log"
	"os"
)

func otherMain(root string) {

	otherExplore(root)

}

func otherExplore(root string) {

	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("error reading dir")

		return
	}

	for _, e := range entries {

		info, err := e.Info()

		if err != nil {
			log.Printf("  (error getting info for %q): %v\n", e.Name(), err)
			continue
		}

		if info.IsDir() {
			otherExplore(root + "/" + info.Name())
		}

	}
}
