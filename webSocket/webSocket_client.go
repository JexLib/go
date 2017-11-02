package webSocket

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type WSClientConfig struct {
	Addr string `flag:"|localhost:8080|http webSocket service address"`
	Path string `flag:"|/ws|http webSocket service path"`
}

type WSClient struct {
	conn *websocket.Conn
}

func NewWSClient(config WSClientConfig) *WSClient {
	u := url.URL{Scheme: "ws", Host: config.Addr, Path: config.Path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	return &WSClient{
		conn: c,
	}
}

type OnClientRecvEvent func(message []byte)

func (wc *WSClient) Start(onevent OnClientRecvEvent) {
	done := make(chan struct{})
	go func() {
		defer wc.conn.Close()
		defer close(done)
		for {
			_, message, err := wc.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			if onevent != nil {
				onevent(message)
			}

		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case t := <-ticker.C:
			err := wc.conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := wc.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			wc.conn.Close()
			return
		}
	}
}
