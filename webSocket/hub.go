// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webSocket

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan jexWsocketBroadcast

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type jexWsocketBroadcast struct {
	msg     []byte
	clients []interface{}
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan jexWsocketBroadcast),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) haveclient(id interface{}, ids []interface{}) bool {
	if len(ids) == 0 {
		return true
	}
	for v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				if h.haveclient(client.id, message.clients) {
					select {
					case client.send <- message.msg:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
