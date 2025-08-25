/*
Copyright 2023 The Numaproj Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"log"

	sideinputsdk "github.com/numaproj/numaflow-go/pkg/sideinput"
)

var counter = 0

// handle is the side input handler function.
func handle(_ context.Context) sideinputsdk.Message {
	// generate message based on even and odd counter
	counter++
	if counter%2 == 0 {
		return sideinputsdk.BroadcastMessage([]byte("even"))
	}
	// BroadcastMessage() is used to broadcast the message with the given value to other side input vertices.
	// val must be converted to []byte.
	return sideinputsdk.BroadcastMessage([]byte("odd"))
}

func main() {
	// Start the side input server.
	err := sideinputsdk.NewSideInputServer(sideinputsdk.RetrieveFunc(handle)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start side input server: ", err)
	}
}
