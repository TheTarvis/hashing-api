package main

import (
	"context"
	"hashing-api/api"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	server := &http.Server{Addr: ":8080", Handler: nil}

	// Sends the SIGTERM and SIGINT calls from the system to the channel to start shutdown process.
	signal.Notify(api.TerminationChannel, syscall.SIGTERM, syscall.SIGINT)
	// This function blocks on api.TerminationChannel
	go shutdown(server)

	//TODO TW: Need to fix issue where url without trailing slash causes request to act like a GET
	http.HandleFunc("/hash/", api.Hash)
	http.HandleFunc("/hash", api.Hash)
	http.HandleFunc("/stats", api.Stats)
	http.HandleFunc("/shutdown", api.Shutdown)

	if err := server.ListenAndServe(); err != nil {
		// http: Server closed gets thrown on server.Shutdown(ctx) and is an expected error.
		if err.Error() != "http: Server closed" {
			log.Printf("HTTP server closed with: %v\n", err)
		}
		// This keeps the sever from shutting down until all the hash jobs have completed.
		api.HashingWaitGroup.Wait()
		log.Printf("HTTP server shut down")
	}

}

func shutdown(server *http.Server) {
	<-api.TerminationChannel // Blocks here until interrupted
	log.Print("Shutdown process initiated\n")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Error while shutting down. %v", err)
	}
}
