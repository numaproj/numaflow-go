package server

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/numaproj/numaflow-go/pkg/info"
	"github.com/stretchr/testify/assert"
)

func Test_getSDKVersion(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "getSDK",
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSDKVersion(); got != tt.want {
				t.Errorf("getSDKVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_information(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	w := httptest.NewRecorder()
	information(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	bs, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	i := &info.ServerInfo{}
	err = json.Unmarshal(bs, i)
	assert.NoError(t, err)
	assert.Equal(t, i.Language, info.Go)
}

func Test_Start(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file, err := os.CreateTemp("/tmp", "test-info.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(file.Name())
		assert.NoError(t, err)
	}()

	go func() {
		if err := Start(ctx, WithSocketAddress(file.Name())); err != nil {
			t.Errorf("Start() error = %v", err)
		}
	}()

	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return net.Dial("unix", file.Name())
	}
	httpClient := &http.Client{
		Transport: httpTransport,
	}
	var resp *http.Response
	i := 3
	for i > 0 {
		resp, err = httpClient.Get("http://unix/ready")
		if err == nil {
			break
		}
		i--
		time.Sleep(1 * time.Second)
	}
	assert.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	infoResp, err := httpClient.Get("http://unix/info")
	assert.NoError(t, err)
	defer func() { _ = infoResp.Body.Close() }()
	assert.Equal(t, http.StatusOK, infoResp.StatusCode)
	bs, err := io.ReadAll(infoResp.Body)
	assert.NoError(t, err)
	si := &info.ServerInfo{}
	err = json.Unmarshal(bs, si)
	assert.NoError(t, err)
	assert.Equal(t, si.Language, info.Go)
	assert.Equal(t, si.Version, "unknown")
}
