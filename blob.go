package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"

	"github.com/charmbracelet/lipgloss/tree"
)

type Config struct {
	maxDepth int
}

type Blob struct {
	name       string
	bytesCount int
	children   map[string]*Blob
}

func newBlob(cfg Config, name string, jsonBlob interface{}) *Blob {
	blobBytes, err := json.Marshal(jsonBlob)
	if err != nil {
		panic(err)
	}

	blob := &Blob{
		name:       name,
		bytesCount: len(blobBytes),
	}

	err = blob.addChildren(cfg, jsonBlob, 0)
	if err != nil {
		panic(err)
	}

	return blob
}

func (b *Blob) addChild(cfg Config, child string, jsonBlob interface{}, currentDepth int) error {
	if b.children == nil {
		b.children = map[string]*Blob{}
	}

	if b.children[child] == nil {
		b.children[child] = &Blob{
			name:       child,
			bytesCount: 0,
		}
	}

	jsonBlobBytes, err := json.Marshal(jsonBlob)
	if err != nil {
		return err
	}
	b.children[child].bytesCount += len(jsonBlobBytes)

	return b.children[child].addChildren(cfg, jsonBlob, currentDepth+1)
}

func (b *Blob) addChildren(cfg Config, jsonBlob interface{}, currentDepth int) error {
	if currentDepth >= cfg.maxDepth {
		return nil
	}

	switch typedElement := jsonBlob.(type) {
	case []interface{}:
		for _, arrayElement := range typedElement {
			b.addChildren(cfg, arrayElement, currentDepth+1)
		}
	case map[string]interface{}:
		for key, childBlob := range typedElement {
			err := b.addChild(cfg, key, childBlob, currentDepth)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Blob) getTree() *tree.Tree {
	blobString := fmt.Sprintf("%s (%s)", b.name, prettyByteSize(b.bytesCount))
	t := tree.Root(blobString)

	childrenArray := []*Blob{}
	for _, child := range b.children {
		childrenArray = append(childrenArray, child)
	}

	sort.Slice(childrenArray, func(i, j int) bool {
		return childrenArray[i].bytesCount > childrenArray[j].bytesCount
	})

	for _, child := range childrenArray {
		t.Child(child.getTree())
	}
	return t
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
