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
	"os"
	"path"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/numaproj/numaflow-go/pkg/mapper"
	"github.com/numaproj/numaflow-go/pkg/sideinput"
)

var sideInputName = "myticker"
var sideInputData []byte
var mu sync.RWMutex

func mapFn(_ context.Context, _ []string, d mapper.Datum) mapper.Messages {
	msg := d.Value()

	// take a read lock before updating sideInputData data to prevent race condition
	mu.RLock()
	siData := sideInputData
	mu.RUnlock()

	if len(siData) > 0 {
		if string(siData) == "even" {
			return mapper.MessagesBuilder().Append(mapper.NewMessage([]byte(string(msg) + "-even-data")))
		} else if string(siData) == "odd" {
			return mapper.MessagesBuilder().Append(mapper.NewMessage([]byte(string(msg) + "-odd-data")))
		}
	}

	return mapper.MessagesBuilder().Append(mapper.MessageToDrop())
}

func main() {
	// Create a new fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Add a path to the watcher
	err = watcher.Add(sideinput.DirPath)
	if err != nil {
		log.Fatal(err)
	}

	// Start a goroutine to listen for events from the watcher
	go fileWatcher(watcher, sideInputName)

	err = mapper.NewServer(mapper.MapperFunc(mapFn)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}

// fileWatcher will watch for any changes in side input file and set data in
// sideInputData global variable.
func fileWatcher(watcher *fsnotify.Watcher, sideInputName string) {
	log.Println("Watching for changes in side input file: ", sideinput.DirPath)
	p := path.Join(sideinput.DirPath, sideInputName)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Println("watcher.Events channel closed")
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create && event.Name == p {
				log.Println("Side input file has been created:", event.Name)
				b, err := os.ReadFile(p)
				if err != nil {
					log.Println("Failed to read side input file: ", err)
				}

				// take a lock before updating sideInputData data to prevent race condition
				mu.Lock()
				sideInputData = b[:]
				mu.Unlock()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				log.Println("watcher.Errors channel closed")
				return
			}
			log.Println("error:", err)
		}
	}
}
