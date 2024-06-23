package main

import (
	"fmt"
	"testing"
)

func TestGetDepthFromPath(t *testing.T) {
	testCases := map[string]int{
		"":           0,
		"/data":      1,
		"/some/path": 2,
	}

	for path, want := range testCases {
		testCase := fmt.Sprintf("getDepthFromPath: %s", path)
		t.Run(testCase, func(t *testing.T) {
			result := getDepthFromPath(path)
			if result != want {
				t.Fatalf("Expected %d, got %d", want, result)
			}
		})
	}
}

func TestParse(t *testing.T) {
	t.Run("Base case", func(t *testing.T) {
		blob := map[string]interface{}{
			"message": "Hello world",
			"nestedBlob": map[string]interface{}{
				"message": "Hello world",
			},
		}

		expectedPaths := []string{"/message", "/nestedBlob/message"}

		pathsToBytes := map[string]int{}
		parse(blob, "", 10, &pathsToBytes)

		for _, path := range expectedPaths {
			_, ok := pathsToBytes[path]
			if !ok {
				t.Fatalf("%s was not picked up", path)
			}
		}
	})

	t.Run("Limit depth", func(t *testing.T) {
		blob := map[string]interface{}{
			"message": "Hello world",
			"nestedBlob": map[string]interface{}{
				"message": "Hello world",
			},
		}

		expectedPaths := []string{"/message"}
		notExpectedPaths := []string{"/nestedBlob/message"}

		pathsToBytes := map[string]int{}
		parse(blob, "", 1, &pathsToBytes)

		for _, path := range expectedPaths {
			_, ok := pathsToBytes[path]
			if !ok {
				t.Fatalf("%s was not picked up", path)
			}
		}

		for _, path := range notExpectedPaths {
			_, ok := pathsToBytes[path]
			if ok {
				t.Fatalf("%s was picked up but not expected", path)
			}
		}
	})

	t.Run("Arrays", func(t *testing.T) {
		blob := map[string]interface{}{
			"message": "Hello world",
			"nestedArray": []interface{}{
				map[string]interface{}{
					"message": "Hello world",
				},
				map[string]interface{}{
					"someOtherMessage": "Goodbye world",
				},
			},
		}

		expectedPaths := []string{"/message", "/nestedArray/message", "/nestedArray/someOtherMessage"}

		pathsToBytes := map[string]int{}
		parse(blob, "", 10, &pathsToBytes)

		for _, path := range expectedPaths {
			_, ok := pathsToBytes[path]
			if !ok {
				t.Fatalf("%s was not picked up", path)
			}
		}
	})
}
