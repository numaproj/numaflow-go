package main

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

var sideInputPath = "/var/numaflow/side-inputs"
var sideInputName = "myticker"

func mapFn(_ context.Context, _ []string, _ mapper.Datum) mapper.Messages {
	// Perform map operation here
	return mapper.MessagesBuilder().Append(mapper.NewMessage([]byte("test_value")))
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
				// Perform some operation here, can update the value in a cache/variable
				log.Println("File contents: ", string(b))
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
