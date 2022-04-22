package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for range jobs {
		data := url.Values{}
		data.Set("password", "asdfasdf")
		responseBody := strings.NewReader(data.Encode())
		//Leverage Go's HTTP Post function to make request
		resp, err := http.Post("http://localhost:8080/hash", "application/x-www-form-urlencoded", responseBody)
		//Handle Error
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()
		//Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Index: %s\n", body)

		atoi, err := strconv.Atoi(string(body))
		if err != nil {
			log.Fatalln(err)
		}
		results <- atoi
	}
}

func main() {

	const numJobs = 500000
	const numWorkers = 1
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		<-results
	}
}
