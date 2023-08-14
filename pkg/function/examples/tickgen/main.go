package main

import (
	"context"
	"encoding/json"
	"time"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

type Data struct {
	Value uint64 `json:"value,omitempty"`
	// only to ensure a desired message size
	Padding []byte `json:"padding,omitempty"`
}

// payload generated by the generator function
// look at newReadMessage function
type payload struct {
	Data      Data
	Createdts int64
}

type ResultPayload struct {
	Value uint64
	Time  string
}

func handle(_ context.Context, keys []string, d functionsdk.Datum) functionsdk.Messages {
	msg := d.Value()
	var p = payload{}
	results := functionsdk.MessagesBuilder()

	var err = json.Unmarshal(msg, &p)
	if err != nil {
		return results
	}

	var v = ResultPayload{
		Value: p.Data.Value,
		Time:  time.Unix(0, p.Createdts).Format(time.RFC3339),
	}
	var r, _ = json.Marshal(v)
	return results.Append(functionsdk.NewMessage(r).WithKeys(keys))
}

func main() {
	server.NewMapServer(context.Background(), functionsdk.MapFunc(handle)).Start(context.Background())
}
