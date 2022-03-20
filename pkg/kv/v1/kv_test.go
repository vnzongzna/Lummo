package kv

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	// test cases
	tests := map[string]struct {
		input  string
		want   string
		status int
	}{
		"simple-fail":   {input: "simple", want: "", status: http.StatusNotFound},
		"exist":         {input: "this", want: `{"value":"that"}`, status: http.StatusOK},
		"alpha-numeric": {input: "alpha-1", want: `{"value":"beta-1"}`, status: http.StatusOK},
	}

	// initialize some values
	kv := Init()
	for k, v := range map[string]string{
		"this":    "that",
		"alpha-1": "beta-1",
	} {
		kv.data[k] = v
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			// need to create request with same parameters as given to chi router
			// pass 'nil' as the third parameter since GET method doesn't expect any query parameter
			req, err := http.NewRequest(http.MethodGet, "/get/{key}", nil)
			if err != nil {
				t.Fatal(err)
			}

			// some chi magic that I found here: https://github.com/go-chi/chi/issues/76#issuecomment-370145140
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("key", tc.input)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

			handler := http.HandlerFunc(kv.Get)
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != tc.status {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.status)
			}

			// Check the returned body is what we expect
			if tc.status == http.StatusOK && rr.Body.String() != tc.want {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.want)
			}
		})
	}
}

func TestSet(t *testing.T) {
	// test cases
	tests := map[string]struct {
		input  string
		status int
	}{
		"simple-fail": {input: `{"sadest-test"}`, status: http.StatusBadRequest},
		"simple-pass": {input: `{"this":"that"}`, status: http.StatusAccepted},
	}

	// initialize nothing
	kv := Init()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPost, "/set", strings.NewReader(tc.input))
			if err != nil {
				t.Fatal(err)
			}

			handler := http.HandlerFunc(kv.Set)
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != tc.status {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.status)
			}
		})
	}

}

func TestSearch(t *testing.T) {
	// test cases
	tests := map[string]struct {
		input  string
		want   string
		status int
	}{
		"simple-fail": {input: "?everythin", want: "", status: http.StatusBadRequest},
		"prefix":      {input: "?prefix=th", want: `{"keys":["this","that"]}`, status: http.StatusOK},
		"suffix":      {input: "?suffix=1", want: `{"keys":["alpha-1","message-1"]}`, status: http.StatusOK},
	}

	// initialize some values
	kv := Init()
	for k, v := range map[string]string{
		"this":      "that",
		"alpha-1":   "beta-2",
		"message-1": "alpha-2",
		"that":      "those",
	} {
		kv.data[k] = v
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/search"+tc.input, nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := http.HandlerFunc(kv.Search)
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != tc.status {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.status)
			}

			// Check the returned body is what we expect
			if tc.status == http.StatusOK && !assert.JSONEq(t, tc.want, rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.want)
			}
		})
	}

}
