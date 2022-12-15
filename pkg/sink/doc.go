// Package sink implements the server code for User Defined Sink in golang.
// Example:
// package main
//
//	 import (
//
//		  "context"
//		  "fmt"
//
//	   sinksdk "github.com/numaproj/numaflow-go/pkg/sink"
//		  "github.com/numaproj/numaflow-go/pkg/sink/server"
//
//	 )
//
//		func handle(ctx context.Context, datumList []sinksdk.Datum) sinksdk.Responses {
//		  result := sinksdk.ResponsesBuilder()
//		  for _, datum := range datumList {
//		    fmt.Println(string(datum.Value()))
//		    result = result.Append(sinksdk.ResponseOK(datum.ID()))
//		  }
//		  return result
//		}
//
//		func main() {
//		  server.New().RegisterSinker(sinksdk.SinkFunc(handle)).Start(context.Background())
//		}
package sink
