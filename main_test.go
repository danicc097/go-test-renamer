package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcessFile(t *testing.T) {
	inputFile := "testdata/input_test.go"
	wantFile := "testdata/input_want.go"

	testdataDir := "testdata"

	tempDir, err := os.MkdirTemp("", "testrenamer")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	if err := copyDir(testdataDir, filepath.Join(tempDir, testdataDir)); err != nil {
		t.Fatalf("Error copying testdata directory: %s", err)
	}

	tempFile := filepath.Join(tempDir, inputFile)

	err = copyFile(inputFile, tempFile)
	if err != nil {
		t.Fatalf("Error copying input file: %s", err)
	}

	err = processFile(tempFile)
	if err != nil {
		t.Fatalf("Error processing file: %s", err)
	}

	got, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Error reading processed file: %s", err)
	}

	want, err := os.ReadFile(wantFile)
	if err != nil {
		t.Fatalf("Error reading file: %s", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0o644)
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	dir, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range dir {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
