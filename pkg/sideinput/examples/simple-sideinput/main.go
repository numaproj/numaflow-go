package main

import (
	"context"
	"time"

	sideinputsdk "github.com/numaproj/numaflow-go/pkg/sideinput"
	"github.com/numaproj/numaflow-go/pkg/sideinput/server"
)

func handle(_ context.Context) sideinputsdk.MessageSI {
	t := time.Now()
	val := "test_value" + string(t.String())
	return sideinputsdk.NewMessageSI([]byte(val))
}
func main() {
	server.NewSideInputServer().RegisterRetriever(sideinputsdk.RetrieveSideInput(handle)).Start(context.Background())
}
