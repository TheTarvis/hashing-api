package api

import (
	"hashing-api/data"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetSumOfRequestTimes(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		{"With no request should return 0", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSumOfRequestTimes(); got != tt.want {
				t.Errorf("GetSumOfRequestTimes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHash(t *testing.T) {
	type args struct {
		method  string
		url     string
		body    string
		headers map[string]string
		hashes  map[int64]string
	}
	type wantedArgs struct {
		statusCode int
		body       string
		headers    map[string]string
	}
	tests := []struct {
		name       string
		args       args
		wantedArgs wantedArgs
	}{
		{"Test empty body returns 400",
			args{
				method: "GET",
				url:    "/hash",
				body:   "",
			},
			wantedArgs{
				statusCode: http.StatusBadRequest,
				body:       "Hash Identifier must be an integer. Value passed: /hash",
			},
		},
		{"Test not allowed method returns method not allowed",
			args{
				method: "DELETE",
				url:    "/hash",
				body:   "",
			},
			wantedArgs{
				statusCode: http.StatusMethodNotAllowed,
				body:       "",
			},
		},
		{"Test get hashed password with id 1 and no hashes saved returns 404",
			args{
				method: "GET",
				url:    "/hash/1",
				body:   "",
			},
			wantedArgs{
				statusCode: http.StatusNotFound,
				body:       "",
			},
		},
		{"Test get hashed password with id 1 and saved hash for that identifier returns value",
			args{
				method: "GET",
				url:    "/hash/1",
				body:   "",
				hashes: map[int64]string{
					1: "asdf",
				},
			},
			wantedArgs{
				statusCode: http.StatusOK,
				body:       "NDAxYjA5ZWFiM2MwMTNkNGNhNTQ5MjJiYjgwMmJlYzhmZDUzMTgxOTJiMGE3NWYyMDFkOGIzNzI3NDI5MDgwZmIzMzc1OTFhYmQzZTQ0NDUzYjk1NDU1NWI3YTA4MTJlMTA4MWMzOWI3NDAyOTNmNzY1ZWFlNzMxZjVhNjVlZDE=",
			},
		},
		{"Test first hash call returns 1",
			args{
				method: "POST",
				url:    "/hash",
				body:   "password=angryMonkey",
				headers: map[string]string{
					"Content-Type": "application/x-www-form-urlencoded",
				},
			},
			wantedArgs{
				statusCode: http.StatusOK,
				body:       "1",
			},
		},
		{"Test empty password returns 400",
			args{
				method: "POST",
				url:    "/hash",
				body:   "password=",
				headers: map[string]string{
					"Content-Type": "application/x-www-form-urlencoded",
				},
			},
			wantedArgs{
				statusCode: http.StatusBadRequest,
				body:       "password must not be empty and must be supplied in the request's form",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			for key, value := range tt.args.headers {
				req.Header.Add(key, value)
			}

			//TODO TW: Need to figure out how to mock this.
			for key, value := range tt.args.hashes {
				data.Get().SavePassword(key, value)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Hash)
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != tt.wantedArgs.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			if rr.Body.String() != tt.wantedArgs.body {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.wantedArgs.body)
			}
		})
	}
}
