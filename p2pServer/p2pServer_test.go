package p2pServer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/kaibling/p2p/libs/Node"
	"github.com/kaibling/p2p/libs/util"
)

func TestGetNodesHandler(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost:54321/getNodes", nil)

	config := &util.Configuration{
		BindingIPAddress: "1.2.3.4",
		BindingPort:      "54321",
		PeerServer:       ""}
	testNode := Node.NewNode("1.2.3.4", "54321")
	server := Newp2pServer(config)

	server.getNodesHandler(res, req)

	var result []Node.Node
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
