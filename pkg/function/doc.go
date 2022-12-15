// Package function implements the server code for User Defined Function in golang.
// Example:
// package main
//
//	 import (
//
//		  "context"
//
//	   functionsdk "github.com/numaproj/numaflow-go/pkg/function"
//		  "github.com/numaproj/numaflow-go/pkg/function/server"
//
//	 )
//
//	 // Simply return the same msg
//
//		func handle(ctx context.Context, key string, data functionsdk.Datum) functionsdk.Messages {
//		  _ = data.EventTime() // Event time is available
//		  _ = data.Watermark() // Watermark is available
//		  return functionsdk.MessagesBuilder().Append(functionsdk.MessageToAll(data.Value()))
//		}
//
//		func main() {
//		  server.New().RegisterMapper(functionsdk.MapFunc(handle)).Start(context.Background())
//		}
package function
