package maradocs

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecodeSuccessBodyNull(t *testing.T) {
	var x struct{ A int }
	if err := decodeSuccessBody([]byte("null"), &x); err != nil {
		t.Fatal(err)
	}
}

func TestTransportAPIException(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(HttpErrorResponse{
			StatusCode: 400,
			ApiErr:     ApiError{Code: 200, Name: "INVALID_SECRET_KEY", Message: "bad"},
		})
	}))
	defer srv.Close()
	tr := newTransport("x", srv.URL, nil)
	var out struct{}
	err := tr.getJSON(t.Context(), "/noop", &out, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *APIException
	if !errors.As(err, &ae) {
		t.Fatalf("want APIException got %T %v", err, err)
	}
	if ae.Details.Code != 200 {
		t.Fatalf("code %d", ae.Details.Code)
	}
}
