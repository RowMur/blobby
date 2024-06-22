package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
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

	parseInterface(blob, "", *maxDepthPtr)
}

func parse(blob map[string]interface{}, path string, maxDepth int) error {
	for key, element := range blob {
		newPath := fmt.Sprintf("%s/%s", path, key)

		elementBytes, err := json.Marshal(element)
		if err != nil {
			return err
		}

		for range getDepthFromPath(newPath) - 1 {
			fmt.Printf("  ")
		}

		fmt.Printf("- %s: %s\n", newPath, prettyByteSize(len(elementBytes)))

		parseInterface(element, newPath, maxDepth)
	}

	return nil
}

func parseInterface(blob interface{}, path string, maxDepth int) {
	if getDepthFromPath(path) >= maxDepth {
		return
	}

	switch typedBlob := blob.(type) {
	case []interface{}:
		for _, arrayBlob := range typedBlob {
			parseInterface(arrayBlob, path, maxDepth)
		}
	case map[string]interface{}:
		parse(typedBlob, path, maxDepth)
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

func getDepthFromPath(path string) int {
	return len(strings.Split(path, "/")) - 1
}
