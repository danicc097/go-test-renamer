package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcessFile(t *testing.T) {
	inputFile := "testdata/input_test.go"
	wantFile := "testdata/input_want.go"

	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatalf("Error reading input file: %s", err)
	}

	got, err := processFile(bytes.NewReader(input))
	if err != nil {
		t.Fatalf("Error processing file: %s", err)
	}

	want, err := os.ReadFile(wantFile)
	if err != nil {
		t.Fatalf("Error reading want file: %s", err)
	}

	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
