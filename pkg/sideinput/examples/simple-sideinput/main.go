package main

import (
	"context"
	"log"
	"time"

	sideinputsdk "github.com/numaproj/numaflow-go/pkg/sideinput"
)

func handle(_ context.Context) sideinputsdk.SideInputMessage {
	t := time.Now()
	val := "test_value" + string(t.String())
	return sideinputsdk.NewSideInputMessage([]byte(val))
}
func main() {
	err := sideinputsdk.NewSideInputServer(sideinputsdk.RetrieveFunc(handle)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start side input server: ", err)
	}
}
