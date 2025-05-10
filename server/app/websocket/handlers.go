package websocket

import (
	"clipboard-share-server/app/person"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	personID, err := person.GetPersonIDFromRequest(w, r)
	if err != nil {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading websocket: %s\n", err.Error())
		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("closing")
		return conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, ""))
	})

	addConnection(personID, conn)
}
