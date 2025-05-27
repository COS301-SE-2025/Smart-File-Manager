package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {

	const root string = "C:/Users/jackb"

	var wg sync.WaitGroup

	start := time.Now()
	wg.Add(1)
	go exploreDir(root, &wg)
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Function execution time: %s\n", elapsed)

	start2 := time.Now()
	otherMain(root)
	fmt.Printf("non conc execution time: %s\n", (time.Since(start2)))

}

//time taken from jacks root directory with concurrency 10 go routine limit:

// max of 10 go routines
var semaphores = make(chan struct{}, 15)

func exploreDir(root string, wg *sync.WaitGroup) {
	defer wg.Done() // ends e routine when this function is done being called. error or not

	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("error reading dir: " + root)

		return
	}

	for _, e := range entries {
		fullPath := filepath.Join(root, e.Name())
		// fmt.Print(fullPath)

		info, err := e.Info()

		if err != nil {
			log.Printf("  (error getting info for %q): %v\n", e.Name(), err)
			continue
		}


		if info.IsDir() {
			wg.Add(1)

			//use an anonymous function as we need the exploredir and ticket grab to be in the same thread
			go func(p string) {
				// defer wg.Done()
				semaphores <- struct{}{} // blocks until it can add a "ticket" to the semaphores

				exploreDir(p, wg)
				<-semaphores //removes a ticket from the semaphores as the func has completed
			}(fullPath)
		}

	}
}
