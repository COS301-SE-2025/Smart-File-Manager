package filesystem

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func SetPreferredCase(w http.ResponseWriter, r *http.Request) {
	// v, found := os.LookupEnv("SFM_SERVER_SECRET")
	// if !found {
	// 	fmt.Println("not ofund")
	// }
	// fmt.Println(v)
	// fmt.Println("v^")
	// preferredCase := r.URL.Query().Get("preferredCase")

	// err := storePreferredCase(preferredCase)
	// if err == nil {
	// 	jsonResponse := []byte(`{"responseMessage": "Saved preferredCase successfully"}`)
	// 	w.Write(jsonResponse)
	// 	return
	// }
	// http.Error(w, ("Failed to write to file: " + err.Error()), http.StatusInternalServerError)

}

func storePreferredCase(preferredCase string) error {
	filePath := filepath.Join("storage", "preferredCase.txt")

	//create file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return err
	}
	defer file.Close()

	data := []byte(preferredCase)

	//write to file
	_, err = file.Write(data)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return err
	}

	return nil
}
