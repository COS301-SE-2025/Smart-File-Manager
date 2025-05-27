package main

import (
	"fmt"
	"log"
	"main/filesystem"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

var (
	composites []*filesystem.Folder
	mu         sync.Mutex
)

func handleComposite(c *filesystem.Folder, request int, path string) {
	mu.Lock()
	defer mu.Unlock()

	switch request {
	case 0: // append
		composites = append(composites, c)
		c.Display(0)
	case 1: // remove
		for i, item := range composites {
			if item.GetPath() == filesystem.ConvertWindowsToWSLPath(path) {
				composites = slices.Delete(composites, i, i+1)
				item = nil
				break
			}
		}
	}

	// Display all stored composites
	fmt.Println("##########################")
	for _, item := range composites {
		fmt.Println("PATH: ", item.GetPath())
	}
	fmt.Println("##########################")
}
func API() {
	fmt.Println("Server started, awaiting requests")
	filesystem.HandleRequests(handleComposite)
}
func main() {
	API()

	// const root string = "C:/Users/jackb"

	// var wg sync.WaitGroup

	// start := time.Now()
	// wg.Add(1)
	// go exploreDir(root, &wg)
	// wg.Wait()

	// elapsed := time.Since(start)
	// fmt.Printf("Function execution time: %s\n", elapsed)

	// start2 := time.Now()
	// otherMain(root)
	// fmt.Printf("non conc execution time: %s\n", (time.Since(start2)))

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

		// fmt.Printf("  Size: %d bytes\n", info.Size())
		// fmt.Printf("  Permissions: %v\n", info.Mode())
		// fmt.Printf("  Modified:    %v\n", info.ModTime())
		// fmt.Printf("  IsDir:       %v\n", info.IsDir())
		// fmt.Println()

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
