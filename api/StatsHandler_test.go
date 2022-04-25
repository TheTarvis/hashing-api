package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStats(t *testing.T) {
	type args struct {
		method string
		url    string
		body   string
		hashes map[int64]string
	}
	type wantedArgs struct {
		statusCode int
		body       *StatsPayload
	}
	tests := []struct {
		name       string
		args       args
		wantedArgs wantedArgs
	}{
		{"Test not allowed method returns method not allowed",
			args{
				method: "DELETE",
				url:    "/stats",
				body:   "",
			},
			wantedArgs{
				statusCode: http.StatusMethodNotAllowed,
				body:       nil,
			},
		},
		{"Test no saved hashes returns 0 total and 0 average",
			args{
				method: "GET",
				url:    "/stats",
				body:   "",
			},
			wantedArgs{
				statusCode: http.StatusOK,
				body: &StatsPayload{
					Total:   0,
					Average: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Stats)
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != tt.wantedArgs.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			wantedBody, _ := json.Marshal(tt.wantedArgs.body)
			if bytes.Compare(rr.Body.Bytes(), wantedBody) == 0 {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), string(wantedBody))
			}
		})
	}
}
