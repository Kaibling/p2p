package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
)


func TestGetNodesHandler(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost:54321/getNodes", nil)

	config := &configuration{
		BindingIPAddress:"1.2.3.4",
		BindingPort:"54321",
		PeerServer:""}
	testNode := newNode("1.2.3.4","54321")
	server := newP2Pserver(config)

	server.getNodesHandler(res, req)

	var result []node
	json.NewDecoder(res.Body).Decode(&result)

	if result[0].IPaddress != testNode.IPaddress {
		t.Errorf("Expected %s, got %s", result[0].IPaddress, testNode.IPaddress)
	}

	if result[0].Port != testNode.Port {
		t.Errorf("Expected %s, got %s", result[0].Port, testNode.Port)
	}

}

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
