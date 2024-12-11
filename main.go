package main

import (
	"log"
	"net/http"
)

// Define a home handler function which writes a byte slice containing
// "Hello from RealQuick" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from RealQsdasuick something"))
}

func main() {
	// Register handler function and corresponding route pattern with the servemux
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Print("starting server on :4000")

	// check via: http://localhost:4000/
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
