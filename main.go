package main

import (
	"fmt"
	"hashing-api/api"
	"net/http"
)

func main() {
	http.HandleFunc("/health", HealthCheck)
	//TODO TW: Need to fix issue where url without trailing slash causes request to act like a GET
	http.HandleFunc("/hash/", api.Hash)
	http.HandleFunc("/hash", api.Hash)

	http.ListenAndServe(":8080", nil)
}

func HealthCheck(writer http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(writer, "Welcome to my website!")
}
