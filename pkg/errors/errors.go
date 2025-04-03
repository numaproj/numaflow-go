package errors

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/numaproj/numaflow-go/pkg/shared"
)

// struct to ensure persist critical error fn is executed only once
// and return error if it has already been executed
type PersistErrorOnce struct {
	done atomic.Bool
	m    sync.Mutex
}
type RuntimeErrorEntry struct {
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

func NewPersistErrorOnce() *PersistErrorOnce {
	return &PersistErrorOnce{}
}

const (
	DEFAULT_RUNTIME_APPLICATION_ERRORS_PATH = "/var/numaflow/runtime/application-errors"
	CURRENT_FILE                            = "current-udf.json"
	INTERNAL_ERROR                          = "Internal error"
)

var containerType = func() string {
	if val, exists := os.LookupEnv(shared.EnvUDContainerType); exists {
		return val
	}
	return "unknown-container"
}()

var persistError = NewPersistErrorOnce()

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
		persistCriticalErrorToFile(errorCode, errorMessage, errorDetails, DEFAULT_RUNTIME_APPLICATION_ERRORS_PATH)
	}
	return nil
}

func persistCriticalErrorToFile(errorCode, errorMessage, errorDetails, dir string) {
	// ModePerm - read/write/execute access permissions to owner, group, and other.
	// Directory may need to be accessible by multiple containers in a containerized environment.
	if dirErr := os.Mkdir(dir, os.ModePerm); dirErr != nil {
		slog.Error("creating directory", "dir", dir, "error", dirErr)
		return
	}
	// Add container to file path
	containerDir := filepath.Join(dir, containerType)
	if dirErr := os.Mkdir(containerDir, os.ModePerm); dirErr != nil {
		slog.Error("creating container directory", "container dir", containerDir, "error", dirErr)
		return
	}

	// Create a current file path
	currentFilePath := filepath.Join(containerDir, CURRENT_FILE)
	f, fileErr := os.Create(currentFilePath)
	if fileErr != nil {
		slog.Error("creating current error log file", "current path", currentFilePath, "error", fileErr)
		return
	}
	defer f.Close()

	if errorCode == "" {
		errorCode = INTERNAL_ERROR
	}

	currentTimestamp := time.Now().Unix()
	runtimeErrorEntry := RuntimeErrorEntry{
		Container: containerType,
		Timestamp: currentTimestamp,
		Code:      errorCode,
		Message:   errorMessage,
		Details:   errorDetails,
	}

	bytesToBeWritten, marshalErr := json.Marshal(runtimeErrorEntry)
	if marshalErr != nil {
		slog.Error("marshalling runtime error entry", "error", marshalErr)
		return
	}
	// Write the error message and details to the current file path
	_, writeErr := f.Write(bytesToBeWritten)
	if writeErr != nil {
		slog.Error("write error to file", "path", currentFilePath, "error", writeErr)
		return
	}

	// Create the final runtime error file path
	fileName := fmt.Sprintf("%d-udf.json", currentTimestamp)
	finalFilePath := filepath.Join(containerDir, fileName)

	// Rename the current file path to final path
	if renameErr := os.Rename(currentFilePath, finalFilePath); renameErr != nil {
		slog.Error("rename path", "current path", currentFilePath, "final path", finalFilePath, "error", renameErr)
		return
	}
}
