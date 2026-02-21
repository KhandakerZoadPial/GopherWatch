package main

import (
	"fmt"
	"gopherwatch/internal/hub"
	"gopherwatch/internal/watcher"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func serveWS(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	clientChannel := make(chan string)
	h.Register <- clientChannel

	defer func() {
		h.Unregister <- clientChannel
	}()

	for {
		message, ok := <-clientChannel
		if !ok {
			break
		}
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			break
		}
	}

}

func main() {
	hub := hub.NewHub()
	go hub.Run()

	watcher := &watcher.FileWatcher{FileName: "test.log"}

	logs, err := watcher.Watch()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Watching test.log... append some text to see it!")

	go func() {
		for line := range logs {
			hub.Broadcast <- line
		}
	}()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("../../ui")))

	fmt.Println("GopherWatch live on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
