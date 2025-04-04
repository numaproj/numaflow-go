package errors

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/numaproj/numaflow-go/pkg/shared"
)

// struct to ensure persist critical error fn is executed only once
// and return error if it has already been executed
type persistErrorOnce struct {
	done atomic.Bool
	m    sync.Mutex
}
type runtimeErrorEntry struct {
	// The name of the container where the error occurred.
	Container string `json:"container"`
	// The timestamp of the error.
	Timestamp int64 `json:"timestamp"`
	// The error code.
	Code string `json:"code"`
	// The error message.
	Message string `json:"message"`
	// Additional details, such as the error stack trace.
	Details string `json:"details"`
}

func newPersistErrorOnce() *persistErrorOnce {
	return &persistErrorOnce{}
}

const (
	runtime_application_errors_path = "/var/numaflow/runtime/application-errors"
	current_file                    = "current-udf.json"
	internal_error                  = "Internal error"
)

var containerType = func() string {
	if val, exists := os.LookupEnv(shared.EnvUDContainerType); exists {
		return val
	}
	return "unknown-container"
}()

var persistError = newPersistErrorOnce()

// PersistCriticalError persists a critical error to an empty dir.
// If the error directory does not exist, it creates it.
// The function will only execute once, regardless of how many times it is called
// Recommended to use this functionality for a critical error in your application
func PersistCriticalError(errorCode, errorMessage, errorDetails string) error {
	if persistError.done.Load() {
		return fmt.Errorf("persist critical error fn executed once already")
	}
	persistError.m.Lock()
	defer persistError.m.Unlock()
	if !persistError.done.Load() {
		defer persistError.done.Store(true)
		if err := persistCriticalErrorToFile(errorCode, errorMessage, errorDetails, runtime_application_errors_path); err != nil {
			fmt.Println("error in persisting critical error: ", err)
		}
	}
	return nil
}

func persistCriticalErrorToFile(errorCode, errorMessage, errorDetails, dir string) error {
	// ModePerm - read/write/execute access permissions to owner, group, and other.
	// Directory may need to be accessible by multiple containers in a containerized environment.
	if dirErr := os.Mkdir(dir, os.ModePerm); dirErr != nil {
		return fmt.Errorf("failed to create directory: %s, error: %w", dir, dirErr)
	}
	// Add container to file path
	containerDir := filepath.Join(dir, containerType)
	if dirErr := os.Mkdir(containerDir, os.ModePerm); dirErr != nil {
		return fmt.Errorf("failed to create container directory: %s, error: %w", containerDir, dirErr)
	}

	// Create a current file path
	currentFilePath := filepath.Join(containerDir, current_file)
	f, fileErr := os.Create(currentFilePath)
	if fileErr != nil {
		return fmt.Errorf("failed to create current error log file: %s, error: %w", currentFilePath, fileErr)
	}
	defer f.Close()

	if errorCode == "" {
		errorCode = internal_error
	}

	currentTimestamp := time.Now().Unix()
	runtimeErrorEntry := runtimeErrorEntry{
		Container: containerType,
		Timestamp: currentTimestamp,
		Code:      errorCode,
		Message:   errorMessage,
		Details:   errorDetails,
	}

	bytesToBeWritten, marshalErr := json.Marshal(runtimeErrorEntry)
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal runtime error entry, error: %w", marshalErr)
	}
	// Write the error message and details to the current file path
	_, writeErr := f.Write(bytesToBeWritten)
	if writeErr != nil {
		return fmt.Errorf("failed to write error to file: %s, error: %w", currentFilePath, writeErr)
	}

	// Create the final runtime error file path
	fileName := fmt.Sprintf("%d-udf.json", currentTimestamp)
	finalFilePath := filepath.Join(containerDir, fileName)

	// Rename the current file path to final path
	if renameErr := os.Rename(currentFilePath, finalFilePath); renameErr != nil {
		return fmt.Errorf("failed to rename file from %s to %s, error: %w", currentFilePath, finalFilePath, renameErr)
	}
	return nil
}
