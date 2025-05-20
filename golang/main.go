package main

import (
	"fmt"
	"log"
	"os"
)

func main(){


	const root string = "./testRootFolder"
	exploreDir(root)
}

func exploreDir(root string){

	

	entries, err := os.ReadDir(root)
	if(err != nil){
		fmt.Println("error reading dir")

		return
	}

	for _, e := range entries{
		fmt.Print(e.Name())
		
        info, err := e.Info()

        if err != nil {
            log.Printf("  (error getting info for %q): %v\n", e.Name(), err)
            continue
        }

        fmt.Printf("  Size: %d bytes\n", info.Size())
        fmt.Printf("  Permissions: %v\n", info.Mode())
        fmt.Printf("  Modified:    %v\n", info.ModTime())
        fmt.Printf("  IsDir:       %v\n", info.IsDir())
        fmt.Println()
		if(info.IsDir()){
			exploreDir(root+"/"+info.Name())
		}

	}
}