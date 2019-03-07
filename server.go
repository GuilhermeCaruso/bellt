package bellt

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

/*
Server - Server startup and standardization structure.
*/
type Server struct {
	Port    string
	Name    string
	Version string
}

// InitServer - Simple initialization of a server to provide a Web or REST
// application.Initializes two routes by default and is responsible for the
// health of the application and for the verification of built routes.
func InitServer(serverInfo Server) {
	var buf bytes.Buffer

	buf.WriteString(":")
	buf.WriteString(serverInfo.Port)
	http.HandleFunc("/health", healthApplication)
	http.Handle("/", verifyBuiltRoutes(redirectBuiltRoute))
	fmt.Printf("%s [%s] was running on port %s\n",
		serverInfo.Name, serverInfo.Version, serverInfo.Port)

	log.Fatal(http.ListenAndServe(buf.String(), nil))
}

// ----------------------------------------------------------------------------
// Server methods
// ----------------------------------------------------------------------------

// Function used in application health routing.
func healthApplication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(HTTPStatusCode["OK"])
	w.Write([]byte(`{"alive": "Server running"}`))
}

// Converts and prepares Handle function for built route analysis
func verifyBuiltRoutes(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
