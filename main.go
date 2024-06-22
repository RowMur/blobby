package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
)

func main() {
	filenamePtr := flag.String("f", "", "the file to parse")
	maxDepthPtr := flag.Int("d", 3, "the maximum depth to parse to")

	flag.Parse()

	filename := *filenamePtr
	if filename == "" {
		panic("Filename not set. Set it with the -f flag.")
	}

	blobBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	var blob interface{}
	err = json.Unmarshal(blobBytes, &blob)
	if err != nil {
		fmt.Println(err)
	}

	parseInterface(blob, 0, *maxDepthPtr)
}

func parse(blob map[string]interface{}, currentDepth int, maxDepth int) error {
	for key, element := range blob {
		elementBytes, err := json.Marshal(element)
		if err != nil {
			return err
		}

		for range currentDepth - 1 {
			fmt.Printf("  ")
		}

		fmt.Printf("- %s: %s\n", key, prettyByteSize(len(elementBytes)))

		parseInterface(element, currentDepth, maxDepth)
	}

	return nil
}

func parseInterface(blob interface{}, currentDepth int, maxDepth int) {
	if currentDepth >= maxDepth {
		return
	}

	switch typedBlob := blob.(type) {
	case []interface{}:
		for _, arrayBlob := range typedBlob {
			parseInterface(arrayBlob, currentDepth, maxDepth)
		}
	case map[string]interface{}:
		parse(typedBlob, currentDepth+1, maxDepth)
	}
}

func prettyByteSize(b int) string {
	bf := float64(b)
	for _, unit := range []string{"", "K", "M", "G"} {
		if math.Abs(bf) < 1024.0 {
			return fmt.Sprintf("%3.1f%sB", bf, unit)
		}
		bf /= 1024.0
	}
	return fmt.Sprintf("%.1fTB", bf)
}
