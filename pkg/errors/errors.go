package errors

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/numaproj/numaflow-go/pkg/shared"
)

const (
	DEFAULT_RUNTIME_APPLICATION_ERRORS_PATH = "/var/numaflow/runtime/application-errors"
	CURRENT_FILE                            = "current-udf.json"
)

var once sync.Once
var containerType = func() string {
	if val, exists := os.LookupEnv(shared.EnvUDContainerType); exists {
		return val
	}
	return "unknown-container"
}()

type RuntimeErrorEntry struct {
	// The name of the container where the error occurred.
	Container string `json:"container"`
	// The timestamp of the error in RFC 3339 format.
	Timestamp string `json:"timestamp"`
	// The error code.
	Code string `json:"code"`
	// The error message.
	Message string `json:"message"`
	// Additional details, such as the error stack trace.
	Details string `json:"details"`
}

// PersistCriticalError persists a critical error to an empty dir.
// If the error directory does not exist, it creates it.
// The function will only execute once, regardless of how many times it is called
// Recommended to use this functionality for a critical error in your application
func PersistCriticalError(errorCode, errorMessage, errorDetails string) error {
	var persistErr error
	once.Do(func() { persistCriticalErrorToFile(&persistErr, errorCode, errorMessage, errorDetails) })
	return persistErr
}

func persistCriticalErrorToFile(err *error, errorCode, errorMessage, errorDetails string) {
	var dir = DEFAULT_RUNTIME_APPLICATION_ERRORS_PATH
	// ModePerm - read/write/execute access permissions to owner, group, and other.
	// Directory may need to be accessible by multiple containers in a containerized environment.
	if dirErr := os.Mkdir(dir, os.ModePerm); dirErr != nil {
		*err = fmt.Errorf("failed to create directory %q: %v", dir, dirErr)
		return
	}
	// Add container to file path
	containerDir := filepath.Join(dir, containerType)
	if dirErr := os.Mkdir(containerDir, os.ModePerm); dirErr != nil {
		*err = fmt.Errorf("failed to create container directory %q: %v", containerDir, dirErr)
		return
	}

	// Create a current file path
	currentFilePath := filepath.Join(containerDir, CURRENT_FILE)
	f, fileErr := os.Create(currentFilePath)
	if fileErr != nil {
		*err = fmt.Errorf("failed to create current error log file %q: %v", currentFilePath, fileErr)
		return
	}
	defer f.Close()

	if errorCode == "" {
		errorCode = "Internal error"
	}

	currentTimestamp := time.Now()
	runtimeErrorEntry := RuntimeErrorEntry{
		Container: containerType,
		Timestamp: currentTimestamp.Format(time.RFC3339),
		Code:      errorCode,
		Message:   errorMessage,
		Details:   errorDetails,
	}

	bytesToBeWritten, marshalErr := json.Marshal(runtimeErrorEntry)
	if marshalErr != nil {
		*err = fmt.Errorf("failed to marshal error entry: %v", marshalErr)
		return
	}
	// Write the error message and details to the current file path
	_, writeErr := f.Write(bytesToBeWritten)
	if writeErr != nil {
		*err = fmt.Errorf("failed to write to error log file %q: %v", currentFilePath, writeErr)
		return
	}

	// Create the final runtime error file path
	fileName := fmt.Sprintf("%d-udf.json", currentTimestamp.Unix())
	finalFilePath := filepath.Join(containerDir, fileName)

	// Rename the current file path to final path
	if renameErr := os.Rename(currentFilePath, finalFilePath); renameErr != nil {
		*err = fmt.Errorf("failed to rename current file path %q to final path %q:  %v", currentFilePath, finalFilePath, writeErr)
		return
	}
}
