package hub

type Hub struct {
	clients map[chan string]bool

	Register   chan chan string
	Unregister chan chan string
	Broadcast  chan string
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[chan string]bool),
		Register:   make(chan chan string),
		Unregister: make(chan chan string),
		Broadcast:  make(chan string),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client)
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				client <- message
			}

		}
	}
}
