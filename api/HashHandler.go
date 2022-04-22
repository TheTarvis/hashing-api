package api

import (
	"fmt"
	"hashing-api/data"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// This sum is of total microseconds for all request to /hash
	sumOfRequestTimes int64
	sumLock           sync.Mutex
	// HashingWaitGroup Used to keep shutdown from killing the server until all tasks for hashing have been completed.
	HashingWaitGroup = &sync.WaitGroup{}
)

func Hash(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	switch request.Method {
	case http.MethodGet:
		getHashHandler(writer, request)
	case http.MethodPost:
		postHashHandler(writer, request)
	}
	duration := time.Since(start)

	go updateSumOfRequestTimes(duration)
}

func updateSumOfRequestTimes(duration time.Duration) {
	sumLock.Lock()
	sumOfRequestTimes += duration.Microseconds()
	sumLock.Unlock()
}

func GetSumOfRequestTimes() int64 {
	return sumOfRequestTimes
}

func getHashHandler(writer http.ResponseWriter, request *http.Request) {
	// Attempting to get the ID from the url by parsing the url path.
	idString := strings.TrimPrefix(request.URL.Path, "/hash/")
	log.Printf("Recieved request to return hash value for id: %s", idString)
	id, err := strconv.Atoi(idString)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprintf(writer, fmt.Sprintf("Hash Identifier must be an integer. Value passed: %s", idString))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error writing response body for bad request.")
			return
		}
		// We return here because we don't want to let the request continue.
		return
	}

	// Attempting to fetch the hashed value from the map
	result, ok := data.Get().GetHashedPassword(int64(id))
	if ok {
		_, err := fmt.Fprintf(writer, fmt.Sprintf("%s", result))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprintf(writer, fmt.Sprint("Error loading password hash."))
			if err != nil {
				return
			}
			return
		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func postHashHandler(writer http.ResponseWriter, request *http.Request) {
	password := request.PostFormValue("password")
	if password == "" {
		writer.WriteHeader(400)
		_, err := fmt.Fprintf(writer, "password must not be empty and must be supplied in the request's form hash.")
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error writing response for Post Hash Handler. %v", err)
			return
		}
		return
	}
	identifier := data.Get().GetNextIdentifier()
	go savePassword(identifier, password)
	_, err := fmt.Fprintf(writer, fmt.Sprintf("%d", identifier))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error writing response for Post Hash Handler. %v", err)
		return
	}
}

func savePassword(id int64, password string) {
	// This allows the shutdown process to wait until the WaitGroup is finished.
	defer HashingWaitGroup.Done()
	HashingWaitGroup.Add(1)
	// Requirement of the coding exercise, would be removed for production.
	time.Sleep(5 * time.Second)
	data.Get().SavePassword(id, password)
}
