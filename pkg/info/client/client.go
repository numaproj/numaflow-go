package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/numaproj/numaflow-go/pkg/info"
)

type InfoClient struct {
	client *http.Client
}

type options struct {
	socketAddress string
	clientTimeout time.Duration
}

type Option func(*options)

func WithSocketAddress(addr string) Option {
	return func(o *options) {
		o.socketAddress = addr
	}
}

func WithClientTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.clientTimeout = timeout
	}
}

// NewInfoClient creates a new info client
func NewInfoClient(opts ...Option) *InfoClient {
	options := &options{
		socketAddress: info.SocketAddress,
		clientTimeout: 3 * time.Second,
	}
	for _, opt := range opts {
		opt(options)
	}

	var httpClient *http.Client
	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return net.Dial("unix", options.socketAddress)
	}
	httpClient = &http.Client{
		Transport: httpTransport,
		Timeout:   options.clientTimeout,
	}
	return &InfoClient{
		client: httpClient,
	}
}

func (c *InfoClient) waitUntilReady(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("failed to wait for ready: %w", ctx.Err())
		default:
			if resp, err := c.client.Get("http://unix/ready"); err == nil {
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
				if resp.StatusCode < 300 {
					return nil
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// GetServerInfo gets the server info
func (c *InfoClient) GetServerInfo(ctx context.Context) (*info.ServerInfo, error) {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := c.waitUntilReady(cctx); err != nil {
		return nil, fmt.Errorf("info server is not ready: %w", err)
	}
	resp, err := c.client.Get("http://unix/info")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http.NewRequestWithContext failed with status, %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll failed: %w", err)
	}
	info := &info.ServerInfo{}
	if err := json.Unmarshal(data, info); err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed: %w", err)
	}
	return info, nil
}
