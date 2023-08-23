package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

var sideInputPath = os.Getenv("SIDEINPUT_PATH")
var sideInputName = os.Getenv("SIDEINPUT_NAME")
var prevValue []byte

func mapFn(_ context.Context, _ []string, _ mapper.Datum) mapper.Messages {
	// Read the side input file.
	p := path.Join(sideInputPath, sideInputName)
	b, err := os.ReadFile(p)
	if err != nil {
		log.Println("Failed to read side input file: ", err)
		return mapper.MessagesBuilder().Append(mapper.MessageToDrop())
	}
	if bytes.Equal(b, prevValue) == false {
		prevValue = b
		return mapper.MessagesBuilder().Append(mapper.NewMessage(b))
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
	err = watcher.Add(sideInputPath)
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

func fileWatcher(watcher *fsnotify.Watcher, sideInputName string) {
	log.Println("Watching for changes in side input file: ", sideInputPath)
	p := path.Join(sideInputPath, sideInputName)
	for {
		select {
		case event, ok := <-watcher.Events:
			log.Println("event:", event)
			if !ok {
				log.Println("watcher.Events channel closed")
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create && event.Name == p {
				log.Println("File has been created:", event.Name)
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
