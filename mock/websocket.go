package mock

import (
	"github.com/gorilla/websocket"
	"sync"
)

type WebSocket struct {
	locker sync.RWMutex
	ws *websocket.Conn
}

func NewWebSocket(ws *websocket.Conn) *WebSocket {
	return &WebSocket{
		ws: ws,
	}
}

func (m *WebSocket) ReadMessage() (int, []byte, error) {
	return m.ws.ReadMessage()
}

func (m *WebSocket) WriteMessage(msgType int, data []byte) error {
	m.locker.Lock()
	defer m.locker.Unlock()
	return  m.ws.WriteMessage(msgType, data)
}

func (m *WebSocket) Close() error {
	return m.ws.Close()
}