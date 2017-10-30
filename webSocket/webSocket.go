package webSocket

import (
	"net/http"

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

func NewWSService() *WSService {

	ws := &WSService{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		Hub: newHub(),
	}
	go ws.Hub.run()
	return ws
}

func (ws *WSService) HandlerHTTP(writer http.ResponseWriter, request *http.Request) error {
	conn, err := ws.Upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return err
	}

	client := &Client{hub: ws.Hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	// clientid := c.Param("timestamp")
	// SSEngine.HandlerHTTP(clientid, c.Response().Writer, c.Request())
	return nil
}

func (ws *WSService) SendMessage(msg []byte) {
	ws.Hub.broadcast <- msg
}
