package main

import (
	"fmt"
	"gopherwatch/internal/watcher"
	"log"
)

func main() {
	watcher := &watcher.FileWatcher{FileName: "test.log"}

	logs, err := watcher.Watch()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Watching test.log... append some text to see it!")

	for line := range logs {
		fmt.Print("New Log: ", line)
	}
}
