package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type p2pserver struct {
	nodeBuffer    *nodeBuffer
	publicIP      string
	configuration *configuration
}
type nodeBuffer struct {
	nodes []node
}

func (nodeBuffer *nodeBuffer) addNode(node node) {
	i := findNodeInArray(nodeBuffer.nodes, node)
	if i == len(nodeBuffer.nodes) {
		log.Println("new element found")
		nodeBuffer.nodes = append(nodeBuffer.nodes, node)
		log.Print("add Node to Buffer: ")
		log.Println(node)
		
	} else {
		log.Println("node already in buffer. skip")
	}
	
}

func (nodeBuffer *nodeBuffer) deleteNode(node node) {

	i := findNodeInArray(nodeBuffer.nodes, node)
	if i == len(nodeBuffer.nodes) {
		log.Println("element not found")
		return
	}
	nodeBuffer.nodes = append(nodeBuffer.nodes[:i], nodeBuffer.nodes[i+1:]...)
	log.Println(nodeBuffer.nodes)

}

func (nodeBuffer *nodeBuffer) toJSON() string {
	jnodes, err := json.Marshal(nodeBuffer.nodes)
	if err != nil {
		log.Fatalln(err)
	}
	return string(jnodes)
}

func newP2Pserver(configuration *configuration) *p2pserver {

	returnP2Pserver := new(p2pserver)
	returnP2Pserver.publicIP = "undef"
	returnP2Pserver.configuration = configuration
	newNode := newNode(configuration.BindingIPAddress, configuration.BindingPort)
	returnP2Pserver.nodeBuffer = new(nodeBuffer)
	returnP2Pserver.nodeBuffer.addNode(newNode)
	return returnP2Pserver

}

func (p2pserver *p2pserver) addNode(ipAddress string, port string) {

	//save locally
	newNode := newNode(ipAddress, port)
	p2pserver.nodeBuffer.addNode(newNode)

	//push to network
	for _, node := range p2pserver.nodeBuffer.nodes {
		//no local connection
		if node.IPaddress == p2pserver.configuration.BindingIPAddress && node.Port == p2pserver.configuration.BindingPort {
			continue
		}
		//todo: no connection to newly published node
		if node.IPaddress == ipAddress && node.Port == port {
			continue
		}

		url := "http://" + node.IPaddress + ":" + node.Port + "/pushNode"
		nodeJSON, err := json.Marshal(node)
		if err != nil {
			log.Println(err)
		}
		postRequest(url,nodeJSON)
	}



}

func (p2pserver *p2pserver) deleteNode(ipAddress string, port string) {
	searchNode := newNode("127.0.0.1", "54321")
	p2pserver.nodeBuffer.deleteNode(searchNode)
	log.Print("Node removed from Buffer: ")

}

func (p2pserver *p2pserver) registerNetwork() {

	connectionString := "http://" + p2pserver.configuration.PeerServer + "/register"
	log.Println("trying to register to " + connectionString)
	localNode := newNode(p2pserver.configuration.BindingIPAddress, p2pserver.configuration.BindingPort)
	bytesRepresentation, err := json.Marshal(localNode)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(connectionString, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result []node
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println("get response from network connect request")
	log.Println(result)
	log.Println("set new node buffer")
	p2pserver.nodeBuffer.nodes = result

}

func (p2pserver *p2pserver) startServer() {
	p2pserver.keepAlive(5)

	if strings.Compare(p2pserver.configuration.PeerServer, "") != 0 {
		log.Println("Connection String " + p2pserver.configuration.PeerServer + " found")
		//connection to server
		p2pserver.registerNetwork()
	} else {
		log.Println("Starting  new network")
	}

	log.Println("Server started on " + p2pserver.configuration.BindingIPAddress + ":" + p2pserver.configuration.BindingPort)

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/getNodes", p2pserver.getNodesHandler)
	http.HandleFunc("/register", p2pserver.registerHandler)
	http.HandleFunc("/pushNode", p2pserver.pushNewNodeInfoHandler)
	
	http.ListenAndServe(":"+p2pserver.configuration.BindingPort, nil)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/ping")
	fmt.Fprintf(w, "OK")

}

func (p2pserver *p2pserver) getNodesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	fmt.Fprintf(w, p2pserver.nodeBuffer.toJSON())
}

func (p2pserver *p2pserver) registerHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("register attempt")
	log.Println(r.Method)
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: %v", err)
		return
	}

	//parse client
	var resa node
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	json.Unmarshal([]byte(buf.String()), &resa)
	log.Println(resa)
	p2pserver.addNode(resa.IPaddress, resa.Port)

	//send to client
	log.Println("send: " + p2pserver.nodeBuffer.toJSON())
	fmt.Fprintf(w, p2pserver.nodeBuffer.toJSON())

}

func (p2pserver *p2pserver) pushNewNodeInfoHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("push new infos ...")
	log.Println(r.Method)
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: %v", err)
		return
	}

	//parse client
	var resa node
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	json.Unmarshal([]byte(buf.String()), &resa)
	log.Println(resa)
	p2pserver.addNode(resa.IPaddress, resa.Port)

	//send to client
	log.Println("send: " + "OK")
	fmt.Fprintf(w, "OK")

}

func (p2pserver *p2pserver) keepAlive(keepAliveTime int) {
	ticker := time.NewTicker(time.Duration(keepAliveTime) * 1000 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for _, node := range p2pserver.nodeBuffer.nodes {
					if node.IPaddress == p2pserver.configuration.BindingIPAddress && node.Port == p2pserver.configuration.BindingPort {
						continue
					}
					url := "http://" + node.IPaddress + ":" + node.Port + "/ping"
					requestData := getRequest(url)
					if requestData == "OK" {
						log.Println("KeepAlive OK with " + url)
						node.setActive()
					} else {
						log.Println("KeepAlive failed with " + url)
						//todo: killt zu schnell
						oldStamp := getHourMinuteSecond(0, 0, -5)
						if node.LastActive.Before(oldStamp) {
							log.Println("node too old")
							p2pserver.nodeBuffer.deleteNode(node)

						}
					}
				}

			}
		}
	}()

	//time.Sleep(3000 * time.Millisecond)
	//ticker.Stop()
	//done <- true
	//fmt.Println("Ticker stopped")
}