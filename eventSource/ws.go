package sse

import (
	"fmt"
	"net/http"
	"time"
)

var (
	wssvr *WSService
)

type WSService struct {
	serveSSE SideEventer
}

func (ws *WSService) HandlerHTTP(clientid interface{}, writer http.ResponseWriter, request *http.Request) {
	ws.serveSSE.HandlerHTTP(clientid, writer, request)
}

func (ws *WSService) SendMessage(msg string, clientid ...interface{}) {
	if len(clientid) > 0 {
		ws.serveSSE.SendEvent(
			&EventOnly{
				CID: clientid,
				Data: &DataEvent{
					Value: msg,
				},
			},
		)
	} else {
		ws.serveSSE.SendEvent(
			&Event{
				Data: &DataEvent{
					Value: msg,
				},
			},
		)
	}
}

// func (w *WSService) ss() {
// 	w.ServeSSE.HandlerHTTP()
// }Writer, c.Request

func Start() *WSService {
	if wssvr != nil {
		return wssvr
	}
	ServeSSE := New(&Config{
		Retry: time.Second * 10,
		Header: map[string]string{
			"Content-Type":  "text/event-stream",
			"Cache-Control": "no-cache",
			"Connection":    "keep-alive",
		},
	})

	ServeSSE.HandlerConnectNotify(func(cid interface{}) {
		// cid is id (client)consumer, when connected and ready
		count := ServeSSE.CountConsumer()
		fmt.Println("client:", cid, " ConnectNotify,count:", count)

		ServeSSE.SendEvent(&EventExcept{
			CID: []interface{}{cid},
			Data: &DataEvent{
				Value: "Connect new user",
			},
		})
	})

	ServeSSE.HandlerDisconnectNotify(func(cid interface{}) {
		// cid is id (client)consumer, when he disconnected from server side event
		count := ServeSSE.CountConsumer()
		ServeSSE.RemoveConsumer(cid)
		fmt.Println("client:", cid, " DisconnectNotify,count:", count)
	})

	ServeSSE.HandlerReconnectNotify(func(rec *Reconnect) {
		// Reconnect struct shows information about client, which reconnected
		// to server side event.
		defer rec.StopRecovery()

		count := ServeSSE.CountConsumer()
		fmt.Println("client ReconnectNotify, count:", count)
	})

	wssvr = &WSService{
		serveSSE: ServeSSE,
	}

	return wssvr
}
