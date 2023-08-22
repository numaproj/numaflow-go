package main

import (
	"context"
	"log"
	"time"

	sideinputsdk "github.com/numaproj/numaflow-go/pkg/sideinput"
)

func handle(_ context.Context) sideinputsdk.MessageSI {
	t := time.Now()
	val := "test_value" + string(t.String())
	return sideinputsdk.NewMessageSI([]byte(val))
}
func main() {
	err := sideinputsdk.NewSideInputServer(sideinputsdk.RetrieverFunc(handle)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start side input server: ", err)
	}
}
