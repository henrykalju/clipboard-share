package routes

import (
	"clipboard-share-server/app/items"
	"net/http"
)

func NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/health", health)

	r.HandleFunc("GET /items", items.GetAllItems)
	r.HandleFunc("POST /items", items.AddItem)

	return r
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
