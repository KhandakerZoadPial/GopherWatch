package hub

type Hub struct {
	clients map[chan string]bool

	Register   chan string
	Unregister chan string
	Broadcast  chan string
}
