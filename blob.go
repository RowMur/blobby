package main

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

var green = lipgloss.ANSIColor(2)
var red = lipgloss.ANSIColor(1)
var blue = lipgloss.ANSIColor(4)

func (b *Blob) getTree() *tree.Tree {
	childrenArray := []*Blob{}
	for _, child := range b.children {
		childrenArray = append(childrenArray, child)
	}

	sort.Slice(childrenArray, func(i, j int) bool {
		return strings.Compare(childrenArray[i].name, childrenArray[j].name) == -1
	})

	blobString := fmt.Sprintf("%s (%s)", b.name, prettyByteSize(b.bytesCount))
	t := tree.Root(blobString).ItemStyleFunc(getItemStyleFunc(childrenArray))

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

func getItemStyleFunc(children []*Blob) tree.StyleFunc {
	lowerBound, upperBound := getOutlierBounds(children)
	return func(_ tree.Children, i int) lipgloss.Style {
		if i < 0 {
			return lipgloss.NewStyle()
		}

		itemBytes := children[i].bytesCount
		if itemBytes > upperBound {
			return lipgloss.NewStyle().Foreground(red)
		} else if itemBytes < lowerBound {
			return lipgloss.NewStyle().Foreground(green)
		}

		return lipgloss.NewStyle().Foreground(blue)
	}
}

func getOutlierBounds(children []*Blob) (int, int) {
	if len(children) == 0 {
		return 0, 0
	}

	values := []float64{}
	for _, child := range children {
		values = append(values, float64(child.bytesCount))
	}

	slices.Sort(values)

	nOfValues := float64(len(values))
	q1Index := int(math.Max((math.Round(nOfValues/4))-1, 0))
	q3Index := int(math.Min(math.Round(3*nOfValues/4)-1, float64(len(values)-1)))

	q1 := values[q1Index]
	q3 := values[q3Index]
	iqr := q3 - q1

	lowerBound := q1 - (1.5 * iqr)
	upperBound := q3 + (1.5 * iqr)

	return int(lowerBound), int(upperBound)
}
