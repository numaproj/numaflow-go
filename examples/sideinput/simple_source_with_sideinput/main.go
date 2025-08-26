package main

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"

	"github.com/numaproj/numaflow-go/pkg/sideinput"
	"github.com/numaproj/numaflow-go/pkg/sourcer"
	"simple_source_with_sideinput/impl"
)

var sideInputName = "myticker"

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
	go impl.FileWatcher(watcher, sideInputName)

	simpleSource := impl.NewSimpleSource()
	err = sourcer.NewServer(simpleSource).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start source server : ", err)
	}
}
