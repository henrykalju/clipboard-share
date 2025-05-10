package routes

import (
	"clipboard-share-server/app/items"
	"clipboard-share-server/app/person"
	"clipboard-share-server/app/webpage"
	"clipboard-share-server/app/websocket"
	"net/http"
)

func NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/health", health)

	r.HandleFunc("GET /login", person.Login)
	r.HandleFunc("POST /register", person.Register)

	r.HandleFunc("GET /items", items.GetAllItems)
	r.HandleFunc("GET /items/{id}", items.GetItemByID)
	r.HandleFunc("POST /items", items.AddItem)

	r.HandleFunc("/ws", websocket.HandleWebsocket)

	r.HandleFunc("/", webpage.HandleWebpage)
	return r
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
