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
	TerminationChannel <- os.Kill
}
