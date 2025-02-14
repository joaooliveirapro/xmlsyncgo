package web

import (
	"fmt"
	"net/http"
)

type WebServer struct {
	DistDir string
}

func (ws *WebServer) Start(port string) {
	// Server index.html
	http.Handle("/", http.FileServer(http.Dir(ws.DistDir)))
	fmt.Println("🟢 server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
