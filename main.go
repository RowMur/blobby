package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	for _, arg := range os.Args {
		if arg == "--help" || arg == "-h" {
			help()
			return
		}
	}

	maxDepthPtr := flag.Int("d", 3, "the maximum depth to parse to")
	rootPtr := flag.String("r", "", "root of the blob to analyse. '.' seperated keys")
	flag.Parse()

	cfg := Config{
		maxDepth: *maxDepthPtr,
	}

	blobBytes := getInput()

	var jsonBlob interface{}
	err := json.Unmarshal(blobBytes, &jsonBlob)
	if err != nil {
		panic(err)
	}

	rootBlob := jsonBlob
	rootPath := strings.Split(*rootPtr, ".")
	for _, key := range rootPath {
		if key == "" {
			continue
		}

		switch typedBlob := rootBlob.(type) {
		case []interface{}:
			index, err := strconv.Atoi(key)
			if err != nil {
				panic(fmt.Sprintf("invalid path at key: %s. Expecting an index to an array", key))
			}

			if index < 0 || index >= len(typedBlob) {
				panic(fmt.Sprintf("invalid index at key: %s. Index out of range for array", key))
			}

			rootBlob = typedBlob[index]
		case map[string]interface{}:
			childBlob, ok := typedBlob[key]
			if !ok {
				panic(fmt.Sprintf("invalid key: %s", key))

			}

			rootBlob = childBlob
		}
	}

	blob := newBlob(cfg, "blob", rootBlob)

	blobTree := blob.getTree()
	fmt.Println(blobTree)
}

func getInput() []byte {
	var blobBytes []byte

	isPiped, err := isPipedInput()
	if err != nil {
		panic(err)
	}

	if isPiped {
		blobBytes, err = io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

		return blobBytes
	}

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

	return fileBlob
}

func isPipedInput() (bool, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}

	return (fi.Mode() & os.ModeCharDevice) == 0, nil
}
