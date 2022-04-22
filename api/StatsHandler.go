package api

import (
	"encoding/json"
	"hashing-api/data"
	"log"
	"net/http"
)

func Stats(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		getStatsHandler(writer, request)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

type StatsPayload struct {
	Total   int64 `json:"total"`
	Average int64 `json:"average"`
}

func getStatsHandler(writer http.ResponseWriter, _ *http.Request) {
	requestCount := data.Get().GetIdentifierCount()
	averageRequestTime := calculateAverageRequestTime(requestCount, GetSumOfRequestTimes())
	response := &StatsPayload{
		Total:   requestCount,
		Average: averageRequestTime,
	}
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Printf("Error while encoding stats payload. %v", err)
		writer.WriteHeader(500)
		return
	}
}

func calculateAverageRequestTime(requestCount int64, totalRequestTime int64) int64 {
	if totalRequestTime == 0 {
		return 0
	}
	return totalRequestTime / requestCount
}
