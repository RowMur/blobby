package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const filename = "blob.json"

func main() {
	blobBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	var blob map[string]interface{}
	err = json.Unmarshal(blobBytes, &blob)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("File: %dB\n", len(blobBytes))
	for key, element := range blob {
		elementBytes, err := json.Marshal(element)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("- %s: %dB\n", key, len(elementBytes))
	}

}
