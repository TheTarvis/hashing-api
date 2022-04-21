package api

import (
	"encoding/base64"
	"fmt"
	"hashing-api/hash"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

//TODO TW: Move the map and count to a data service to clean this file up.
var hashCount = 0
var m sync.Mutex
var sm sync.Map

func Hash(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		hashGetHandler(writer, request)
	case http.MethodPost:
		hashPostHandler(writer, request)
	}
}

func hashGetHandler(writer http.ResponseWriter, request *http.Request) {
	// Attempting to get the ID from the url by parsing the url path.
	idString := strings.TrimPrefix(request.URL.Path, "/hash/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, fmt.Sprintf("Hash Identifier must be an integer. Value passed: %s", idString))
		return
	}

	// Attempting to fetch the hashed value from the map
	result, ok := sm.Load(id)
	if ok {
		_, err := fmt.Fprintf(writer, fmt.Sprintf("%s", result))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(writer, fmt.Sprint("Error loading password hash."))
			return
		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func hashPostHandler(writer http.ResponseWriter, request *http.Request) {
	// Making sure this  goroutine gets the correct counter without a collision
	var w sync.WaitGroup
	w.Add(1)
	go increment(&w, &m)
	w.Wait()

	password := request.PostFormValue("password")
	if password == "" {
		writer.WriteHeader(400)
		fmt.Fprintf(writer, "password must not be empty and must be supplied in the request's form data.")
		return
	}
	go hashAndStorePassword(hashCount, password)
	fmt.Fprintf(writer, fmt.Sprintf("%d", hashCount))
}

func increment(wg *sync.WaitGroup, m *sync.Mutex) {
	m.Lock()
	hashCount = hashCount + 1
	m.Unlock()
	wg.Done()
}

func hashAndStorePassword(id int, password string) {
	time.Sleep(5 * time.Second)
	value := hash.Sha512(password)
	sm.Store(id, base64.StdEncoding.EncodeToString([]byte(value)))
}
