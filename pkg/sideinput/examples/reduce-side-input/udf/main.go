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
	"strconv"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/numaproj/numaflow-go/pkg/reducer"
	"github.com/numaproj/numaflow-go/pkg/sideinput"
)

// SumReducerCreator implements the reducer.ReducerCreator interface which creates a reducer
type SumReducerCreator struct {
}

func (s *SumReducerCreator) Create() reducer.Reducer {
	return &Sum{}
}

// Sum is a reducer that sum up the values for the given keys
type Sum struct {
	sum int
}

var sideInputName = "myticker"
var sideInputData []byte
var mu sync.RWMutex

func (s *Sum) Reduce(ctx context.Context, keys []string, reduceCh <-chan reducer.Datum, md reducer.Metadata) reducer.Messages {
	mu.RLock()
	siData := sideInputData
	mu.RUnlock()

	// sum up the values
	for d := range reduceCh {
		value, err := strconv.Atoi(string(d.Value()))
		if err != nil {
			log.Printf("unable to convert the value to integer: %v\n", err)
			continue
		}

		s.sum += value
	}

	if s.sum > 0 && len(siData) > 0 {
		msg := []byte("reduce-side-input")
		return reducer.MessagesBuilder().Append(reducer.NewMessage(msg))
	} else {
		return reducer.MessagesBuilder().Append(reducer.MessageToDrop())
	}
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

	err = reducer.NewServer(&SumReducerCreator{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start reducer function server: ", err)
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
