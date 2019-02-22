package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/GuilhermeCaruso/bellt/utils"
)

/*
InitServer - Simple server initialization
*/
func InitServer(serverInfo Server) {
	var buf bytes.Buffer

	buf.WriteString(":")
	buf.WriteString(serverInfo.Port)
	http.HandleFunc("/health", healthApplication)

	fmt.Printf("%s [%s] was running on port %s", serverInfo.Name, serverInfo.Version, serverInfo.Port)

	log.Fatal(http.ListenAndServe(buf.String(), nil))
}

//Simple function to application health check
func healthApplication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(utils.HTTPStatusCode["OK"])
	w.Write([]byte(`{"alive": "Server running"}`))
}
