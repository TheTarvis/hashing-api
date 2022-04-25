package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//TODO TW: Need to figure out how to mock channels
func TestShutdown(t *testing.T) {
	type args struct {
		method  string
		url     string
		body    string
		headers map[string]string
	}
	type wantedArgs struct {
		statusCode int
	}
	tests := []struct {
		name       string
		args       args
		wantedArgs wantedArgs
	}{
		{"Test shutdown add signal to terminiation channel",
			args{
				method:  "GET",
				url:     "/shutdown",
				body:    "",
				headers: nil,
			},
			wantedArgs{
				statusCode: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				req, _ := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
				for key, value := range tt.args.headers {
					req.Header.Add(key, value)
				}

				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(Shutdown)
				handler.ServeHTTP(rr, req)

				// Check the status code is what we expect.
				if status := rr.Code; status != tt.wantedArgs.statusCode {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, http.StatusOK)
				}
			})
		})
	}
}
