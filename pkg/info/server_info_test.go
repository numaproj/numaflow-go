package info

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func Test_getSDKVersion(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "getSDK",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSDKVersion(); got != tt.want {
				t.Errorf("GetSDKVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_WaitUntilReady(t *testing.T) {
	serverInfoFile, err := os.CreateTemp("/tmp", "server-info")
	assert.NoError(t, err)
	defer os.Remove(serverInfoFile.Name())
	err = os.WriteFile(serverInfoFile.Name(), []byte("test"), 0644)
	assert.NoError(t, err)

	t.Run("test timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		err := WaitUntilReady(ctx, WithServerInfoFilePath("/tmp/not-exist"))
		assert.True(t, errors.Is(err, context.DeadlineExceeded))
	})

	t.Run("test success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		err = WaitUntilReady(ctx, WithServerInfoFilePath(serverInfoFile.Name()))
		assert.NoError(t, err)
	})
}

func Test_Read_Write(t *testing.T) {
	filepath := os.TempDir() + "/server-info"
	defer os.Remove(filepath)
	info := &ServerInfo{
		Protocol:               TCP,
		Language:               Java,
		MinimumNumaflowVersion: MinimumNumaflowVersion,
		Version:                "11",
		Metadata:               map[string]string{"key1": "value1", "key2": "value2"},
	}
	err := Write(info, WithServerInfoFilePath(filepath))
	assert.NoError(t, err)
	got, err := Read(WithServerInfoFilePath("/tmp/not-exist"))
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
	assert.Nil(t, got)
	got, err = Read(WithServerInfoFilePath(filepath))
	assert.NoError(t, err)
	assert.Equal(t, info, got)
}
