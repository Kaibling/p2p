package peerserver

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)
func TestPingHandler(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost:54321/ping", nil)

	pingHandler(res, req)

	content, _ := ioutil.ReadAll(res.Body)
	expected := "OK"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s", expected, string(content))
	}
}
