package api

import (
	"log"
	"net/http"
	"os"
)

var TerminationChannel = make(chan os.Signal)

func Shutdown(_ http.ResponseWriter, _ *http.Request) {
	log.Println("Starting shutdown process.")
	// Puts an os.Kill signal on the termination channel to start shutdown
	// put it in a go routine so the test would not wait, this is probably way wrong.
	go func() { TerminationChannel <- os.Kill }()
}
