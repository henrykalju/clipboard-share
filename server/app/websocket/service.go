package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
	"slices"
)

const HISTORY_UPDATED_MESSAGE = "HISTORY"

var (
	clients      = make(map[int32][]*websocket.Conn)
	clientsMutex = sync.Mutex{}
)

func addConnection(person int32, conn *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	_, ok := clients[person]
	if !ok {
		clients[person] = make([]*websocket.Conn, 0, 1)
	}

	clients[person] = append(clients[person], conn)
	go handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
	for {
		_, _, err := conn.NextReader()
		if err != nil {
			break
		}
	}

	clientsMutex.Lock()
	for key, value := range clients {
		for i := len(value) - 1; i >= 0; i-- {
			if conn == value[i] {
				clients[key] = slices.Delete(clients[key], i, i+1)
			}
		}
	}
	clientsMutex.Unlock()
}

func NotifyConnectionsOfHistoryUpdate(person int32) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	conns, ok := clients[person]
	if !ok {
		return
	}

	for _, conn := range conns {
		conn.WriteMessage(websocket.TextMessage, []byte(HISTORY_UPDATED_MESSAGE))
	}
}
