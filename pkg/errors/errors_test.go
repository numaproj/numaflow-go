package errors

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/numaproj/numaflow-go/pkg/internal/shared"
)

func TestPersistCriticalErrorToFileWritesErrorDetails(t *testing.T) {
	dir := "testdirone"
	defer os.RemoveAll(dir)

	errorCode, errorMessage, errorDetails := "testCode", "testMessage", "testDetails"

	if err := persistCriticalErrorToFile(errorCode, errorMessage, errorDetails, dir); err != nil {
		t.Fatalf("failed to persist error: %v", err)
	}

	containerDir := filepath.Join(dir, shared.ContainerType)
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
	dir := "testdirtwo"
	defer os.RemoveAll(dir)

	if err := persistCriticalErrorToFile("", "message", "details", dir); err != nil {
		t.Fatalf("failed to persist error: %v", err)
	}

	containerDir := filepath.Join(dir, shared.ContainerType)
	fileName := fmt.Sprintf("%d-udf.json", time.Now().Unix())
	finalFilePath := filepath.Join(containerDir, fileName)

	if _, err := os.Stat(finalFilePath); os.IsNotExist(err) {
		t.Errorf("expected file %s to be created, but it does not exist", finalFilePath)
	}

	f, err := os.Open(finalFilePath)
	if err != nil {
		t.Fatalf("failed to open file %s: %v", finalFilePath, err)
	}
	defer f.Close()

	var entry runtimeErrorEntry
	if err := json.NewDecoder(f).Decode(&entry); err != nil {
		t.Fatalf("failed to decode JSON from file %s: %v", finalFilePath, err)
	}

	if entry.Code != internal_error {
		t.Errorf("expected error code to be %s, but got %s", internal_error, entry.Code)
	}
}

func TestPersistCriticalErrorToFileCreatesDirectory(t *testing.T) {
	dir := "testdirthree"
	defer os.RemoveAll(dir)

	persistCriticalErrorToFile("code", "message", "details", dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("expected directory %s to be created, but it does not exist", dir)
	}
}

func TestPersistCriticalErrorAllReturnError(t *testing.T) {
	errorCode, errorMessage, errorDetails := "testCode", "testMessage", "testDetails"

	persistError.done.Store(true)
	defer persistError.done.Store(false)

	var wg sync.WaitGroup
	numGoroutines := 10
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errors <- PersistCriticalError(errorCode, errorMessage, errorDetails)
		}()
	}

	wg.Wait()
	close(errors)

	failCount := 0
	for err := range errors {
		if err != nil && strings.Contains(err.Error(), "persist critical error fn executed once already") {
			failCount++
		}
	}

	if failCount != numGoroutines {
		t.Errorf("expected all %d goroutines to fail, but only %d failed", numGoroutines, failCount)
	}
}
