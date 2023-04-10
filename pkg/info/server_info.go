package info

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func GetSDKVersion() string {
	version := "unknown"
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return version
	}
	for _, d := range info.Deps {
		if strings.HasSuffix(d.Path, "/numaflow-go") {
			version = d.Version
			break
		}
	}
	return version
}

// Write writes the server info to a file
func Write(svrInfo *ServerInfo, opts ...Option) error {
	b, err := json.Marshal(svrInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal server info: %w", err)
	}
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	if err := os.Remove(options.svrInfoFilePath); !os.IsNotExist(err) && err != nil {
		return fmt.Errorf("failed to remove server-info file: %w", err)
	}
	if err := os.WriteFile(options.svrInfoFilePath, b, 0644); err != nil {
		return fmt.Errorf("failed to write server-info file: %w", err)
	}
	return nil
}

// WaitUntilReady waits until the server info is ready
func WaitUntilReady(ctx context.Context, opts ...Option) error {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
existence:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if _, err := os.Stat(options.svrInfoFilePath); err != nil {
				log.Printf("Server info file %s is not ready...", options.svrInfoFilePath)
				time.Sleep(1 * time.Second)
				continue
			} else {
				break existence
			}
		}
	}

	// Also check if the file is ready to read
	// It takes some time for the server to write the server info file
	// TODO: use a better way to wait for the file to be ready
	b, err := os.ReadFile(options.svrInfoFilePath)
	if err != nil {
		return err
	}
	for len(b) == 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(100 * time.Millisecond)
			b, err = os.ReadFile(options.svrInfoFilePath)
			if err != nil {
				return err
			}
			if len(b) > 0 {
				return nil
			}
		}
	}
	return nil
}

// Read reads the server info from a file
func Read(opts ...Option) (*ServerInfo, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	// It takes some time for the server to write the server info file
	// TODO: use a better way to wait for the file to be ready
	retry := 0
	b, err := os.ReadFile(options.svrInfoFilePath)
	for len(b) == 0 && err == nil && retry < 10 {
		time.Sleep(100 * time.Millisecond)
		b, err = os.ReadFile(options.svrInfoFilePath)
		retry++
	}
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return nil, fmt.Errorf("server info file is empty")
	}
	info := &ServerInfo{}
	if err := json.Unmarshal(b, info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal server info: %w", err)
	}
	return info, nil
}
