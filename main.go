package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	maxDepthPtr := flag.Int("d", 3, "the maximum depth to parse to")
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

	blob := newBlob(cfg, "blob", jsonBlob)

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
