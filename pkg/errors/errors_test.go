package errors

import (
	"encoding/json"
	"fmt"

	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPersistCriticalErrorToFileWritesErrorDetails(t *testing.T) {
	dir := "testdir"
	defer os.RemoveAll(dir)

	errorCode := "testCode"
	errorMessage := "testMessage"
	errorDetails := "testDetails"

	persistCriticalErrorToFile(errorCode, errorMessage, errorDetails, dir)

	containerDir := filepath.Join(dir, containerType)
	files, err := os.ReadDir(containerDir)
	if err != nil {
		t.Fatalf("failed to read container directory: %v", err)
	}

	if len(files) == 0 {
		t.Fatalf("expected a file to be created in the container directory, but none was found")
	}

	finalFilePath := filepath.Join(containerDir, files[0].Name())
	content, err := os.ReadFile(finalFilePath)
	if err != nil {
		t.Fatalf("failed to read the error file: %v", err)
	}

	if !strings.Contains(string(content), errorCode) || !strings.Contains(string(content), errorMessage) || !strings.Contains(string(content), errorDetails) {
		t.Errorf("error file content does not match expected values")
	}
}

func TestPersistCriticalErrorToFileUsesDefaultErrorCode(t *testing.T) {
	dir := "testdir"
	defer os.RemoveAll(dir)

	persistCriticalErrorToFile("", "message", "details", dir)

	containerDir := filepath.Join(dir, containerType)
	fileName := fmt.Sprintf("%d-udf.json", time.Now().Unix())
	finalFilePath := filepath.Join(containerDir, fileName)

	if _, err := os.Stat(finalFilePath); os.IsNotExist(err) {
		t.Errorf("expected file %s to be created, but it does not exist", finalFilePath)
	}

	// Additional check to ensure the default error code is used
	f, err := os.Open(finalFilePath)
	if err != nil {
		t.Fatalf("failed to open file %s: %v", finalFilePath, err)
	}
	defer f.Close()

	var entry RuntimeErrorEntry
	if err := json.NewDecoder(f).Decode(&entry); err != nil {
		t.Fatalf("failed to decode JSON from file %s: %v", finalFilePath, err)
	}

	if entry.Code != INTERNAL_ERROR {
		t.Errorf("expected error code to be %s, but got %s", INTERNAL_ERROR, entry.Code)
	}
}

func TestPersistCriticalErrorToFileCreatesDirectory(t *testing.T) {
	dir := "testdir"
	defer os.RemoveAll(dir)

	persistCriticalErrorToFile("code", "message", "details", dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("expected directory %s to be created, but it does not exist", dir)
	}
}
