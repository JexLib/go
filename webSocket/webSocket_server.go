package webSocket

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

/** 例子
*<script>
    try {

        var timestamp = new Date().getTime();
        var sock = new WebSocket("ws://{{.Domain}}/ws");
        //sock.binaryType = 'blob'; // can set it to 'blob' or 'arraybuffer
        console.log("Websocket - status: " + sock.readyState);
        sock.onopen = function(m) {
            console.log("CONNECTION opened..." + this.readyState);
        }
        sock.onmessage = function(m) {
            $('#chatbox').append('<p>' + m.data + '</p>');
        }
        sock.onerror = function(m) {
            console.log("Error occured sending..." + m.data);
        }
        sock.onclose = function(m) {
            console.log("Disconnected - status " + this.readyState);
        }
    } catch (exception) {
        console.log(exception);
    }
</script>
*/
// type Config struct {
// 	ReadBufferSize  int
// 	WriteBufferSize int
// 	Newline         string
// 	space           string
// }

type WSService struct {
	Upgrader websocket.Upgrader
	Hub      *Hub
}

type OnClientEvent func(client *Client)

func NewWSService() *WSService {

	ws := &WSService{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			EnableCompression: true,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Hub: newHub(),
	}
	go ws.Hub.run()
	return ws
}

func (ws *WSService) SetClientRegeditEvent(onClientEvent OnClientEvent) {
	ws.Hub.onClientRegedit = onClientEvent
}

func (ws *WSService) SetClientUnRegeditEvent(onClientEvent OnClientEvent) {
	ws.Hub.onClientUnRegedit = onClientEvent
}

func (ws *WSService) HandlerHTTP(writer http.ResponseWriter, request *http.Request, clientid ...interface{}) error {
	conn, err := ws.Upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return err
	}

	client := &Client{hub: ws.Hub, conn: conn, send: make(chan []byte, 256), UUID: uuid.New()}
	if len(clientid) > 0 {
		client.ID = clientid[0]
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	// clientid := c.Param("timestamp")
	// SSEngine.HandlerHTTP(clientid, c.Response().Writer, c.Request())

	return nil
}

func (ws *WSService) SendMessage(msg []byte, clientid ...interface{}) {
	ws.Hub.broadcast <- jexWsocketBroadcast{msg: msg, clients: clientid}
}
