package main

import (
	"os"
	"strconv"
	"testing"
)

var lineFiles []*LineFile

func TestMain(m *testing.M) {
	// The small test file has 5 lines, 10 lines per index pages will result in one index page
	lineFiles = append(lineFiles, NewLineFile("./text_file_small1.txt", 10))
	// The small test file has 5 lines, 2 lines per index pages will result in more than one index page
	lineFiles = append(lineFiles, NewLineFile("./text_file_small2.txt", 2))
	// The small test file has 1,000,000 lines
	lineFiles = append(lineFiles, NewLineFile("./text_file_1M_line.txt", 1000))
	os.Exit(m.Run())
}

// TestLineFileBuildIndex tests index building on all LineFile objects
func TestLineFileBuildIndex(t *testing.T) {
	for i := 0; i < len(lineFiles); i++ {
		lineFile := lineFiles[i]
		lineFile.BuildIndex()
		if !lineFile.indexCompleted {
			t.Errorf("indexCompleted is not correct: got %v, expected: %v", lineFile.indexCompleted, true)
		}
	}
}

// TestLineFileLineGet tests "get line by line number" on the first two LineFile objects
// - one with one index page, the other with multiple index pages
func TestLineFileLineGet(t *testing.T) {
	expectedLines := []string{"a", "b", "c", "d"}
	for f := 0; f < 2; f++ {
		for i := 1; i < 5; i++ {
			lineFile := lineFiles[f]
			status, line := lineFile.GetLine(i)
			if status != 200 || line != expectedLines[i-1] {
				t.Errorf("Get line at %d: got %s, expected: %s", i, line, expectedLines[i])
			}
		}
	}
}

// TestLineFileErrCases tests "get line by invalid line number" on the first two LineFile objects
// - one with one index page, the other with multiple index pages
func TestLineFile1MLineFile(t *testing.T) {
	lineFile := lineFiles[2]
	for i := 1; i <= 10; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			i := GetRandom(1000000)
			status, _ := lineFile.GetLine(i)
			if status != 200 {
				t.Errorf("Get line at %d: got status %d", i, status)
			}
		})
	}
}

// TestLineFileErrCases tests "get line by invalid line number" on the first two LineFile objects
// - one with one index page, the other with multiple index pages
func TestLineFileErrCases(t *testing.T) {
	for f := 0; f < 2; f++ {
		lineFile := lineFiles[f]
		status, _ := lineFile.GetLine(0)
		if status != 404 {
			t.Errorf("Get line at 0: status - got %d, expected: %d", status, 404)
		}
		// The first two LineFile objects has 5 lines so line 6 should return 404 - not found.
		status2, _ := lineFile.GetLine(6)
		if status2 != 404 {
			t.Errorf("Get line at 0: status - got %d, expected: %d", status, 404)
		}
	}
}
