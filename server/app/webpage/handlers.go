package webpage

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed index.html
var indexHtml []byte

func HandleWebpage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write(indexHtml)
	if err != nil {
		fmt.Printf("Error responding with index.html: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
