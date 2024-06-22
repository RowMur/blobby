package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	filenamePtr := flag.String("f", "", "the file to parse")
	flag.Parse()

	filename := *filenamePtr
	if filename == "" {
		panic("Filename not set. Set it with the -f flag.")
	}

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
