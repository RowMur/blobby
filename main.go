package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
)

func main() {
	maxDepthPtr := flag.Int("d", 3, "the maximum depth to parse to")
	flag.Parse()

	var blobBytes []byte

	isPiped, err := isPipedInput()
	if err != nil {
		panic(err)
	}

	if isPiped {
		_, err = fmt.Scanf("%s", &blobBytes)
		if err != nil {
			panic(err)
		}
	} else {
		args := os.Args[1:]
		if len(args) == 0 {
			panic("Filename not set.")
		}

		filename := args[len(args)-1]
		if filename == "" {
			panic("Filename not set.")
		}

		fileBlob, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		blobBytes = fileBlob
	}

	var blob interface{}
	err = json.Unmarshal(blobBytes, &blob)
	if err != nil {
		panic(err)
	}

	blobPathsToBytes := &map[string]int{}
	parseInterface(blob, "", *maxDepthPtr, blobPathsToBytes)

	pathsArray := []string{}
	for path := range *blobPathsToBytes {
		pathsArray = append(pathsArray, path)
	}
	slices.Sort(pathsArray)

	for _, path := range pathsArray {
		splitPath := strings.Split(path, "/")
		for range len(splitPath) - 2 {
			fmt.Printf("  ")
		}
		fmt.Printf("- %s: %s\n", splitPath[len(splitPath)-1], prettyByteSize((*blobPathsToBytes)[path]))
	}
}

func parse(blob map[string]interface{}, path string, maxDepth int, pathsToBytes *map[string]int) error {
	for key, element := range blob {
		elementBytes, err := json.Marshal(element)
		if err != nil {
			return err
		}

		newPath := fmt.Sprintf("%s/%s", path, key)
		_, ok := (*pathsToBytes)[newPath]
		if !ok {
			(*pathsToBytes)[newPath] = len(elementBytes)
		} else {
			(*pathsToBytes)[newPath] += len(elementBytes)
		}

		parseInterface(element, newPath, maxDepth, pathsToBytes)
	}

	return nil
}

func parseInterface(blob interface{}, path string, maxDepth int, pathsToBytes *map[string]int) {
	if getDepthFromPath(path) >= maxDepth {
		return
	}

	switch typedBlob := blob.(type) {
	case []interface{}:
		for _, arrayBlob := range typedBlob {
			parseInterface(arrayBlob, path, maxDepth, pathsToBytes)
		}
	case map[string]interface{}:
		parse(typedBlob, path, maxDepth, pathsToBytes)
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

func isPipedInput() (bool, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}

	return (fi.Mode() & os.ModeCharDevice) == 0, nil
}
