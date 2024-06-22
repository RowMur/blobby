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

	parse(blob, 0)
}

func parse(blob map[string]interface{}, currentDepth int) error {
	for key, element := range blob {
		elementBytes, err := json.Marshal(element)
		if err != nil {
			return err
		}

		for range currentDepth {
			fmt.Printf("  ")
		}

		fmt.Printf("- %s: %dB\n", key, len(elementBytes))

		parseInterface(element, currentDepth)
	}

	return nil
}

func parseInterface(blob interface{}, currentDepth int) {
	switch typedBlob := blob.(type) {
	case []interface{}:
		for _, arrayBlob := range typedBlob {
			parseInterface(arrayBlob, currentDepth)
		}
	case map[string]interface{}:
		parse(typedBlob, currentDepth+1)
	}
}
