package p2pServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/Kaibling/p2p/libs/Node"
	"github.com/Kaibling/p2p/libs/util"
)

type Payload interface {
	getVersion() string
	saveData(interface {})
	getData() interface {}
	toJSON() string
}

type p2pserver struct {
	nodeBuffer    *nodeBuffer
	publicIP      string
	configuration *util.Configuration
	payload			*Payload
}
type nodeBuffer struct {
	nodes []Node.Node
}

type configCheck struct {
	KeepAlive time.Duration
	NetworkName string
}

func (nodeBuffer *nodeBuffer) addNode(node Node.Node) {
	i := util.FindNodeInArray(nodeBuffer.nodes, node)
	if i == len(nodeBuffer.nodes) {
		log.Println("new element found")
		nodeBuffer.nodes = append(nodeBuffer.nodes, node)
		log.Print("add Node to Buffer: ")
		log.Println(node)

	} else {
		log.Println("node already in buffer. skip")
	}

}

func (nodeBuffer *nodeBuffer) deleteNode(node Node.Node) {

	i := util.FindNodeInArray(nodeBuffer.nodes, node)
	if i == len(nodeBuffer.nodes) {
		log.Println("element not found")
		return
	}
	nodeBuffer.nodes = append(nodeBuffer.nodes[:i], nodeBuffer.nodes[i+1:]...)
	log.Println(node)
	log.Println("node deleted")

}

func (nodeBuffer *nodeBuffer) toJSON() string {
	jnodes, err := json.Marshal(nodeBuffer.nodes)
	if err != nil {
		log.Fatalln(err)
	}
	return string(jnodes)
}


func Newp2pServer(configuration *util.Configuration) *p2pserver {

	returnP2Pserver := new(p2pserver)
	returnP2Pserver.publicIP = "undef"
	returnP2Pserver.configuration = configuration
	newNode := Node.NewNode(configuration.BindingIPAddress, configuration.BindingPort)
	returnP2Pserver.nodeBuffer = new(nodeBuffer)
	returnP2Pserver.nodeBuffer.addNode(newNode)
	return returnP2Pserver

}

func (p2pserver *p2pserver) pushNode(ipAddress string, port string) {
	newNode := Node.NewNode(ipAddress, port)
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
		nodeJSON, err := json.Marshal(newNode)
		if err != nil {
			log.Println(err)
		}
		util.PostRequest(url, nodeJSON)
	}

}

func (p2pserver *p2pserver) getPayload() *Payload {
	return p2pserver.payload
}

func (p2pserver *p2pserver) addNode(ipAddress string, port string) {

	//save locally
	newNode := Node.NewNode(ipAddress, port)
	p2pserver.nodeBuffer.addNode(newNode)

}

func (p2pserver *p2pserver) deleteNode(ipAddress string, port string) {
	searchNode := Node.NewNode(ipAddress, port)
	p2pserver.nodeBuffer.deleteNode(searchNode)
}

func (p2pserver *p2pserver) AddPayload(payload *Payload) {
	p2pserver.payload = payload

}

func (p2pserver *p2pserver) registerNetwork() {

	//send own node data to server
	connectionString := "http://" + p2pserver.configuration.PeerServer + "/register"
	log.Println("trying to register to " + connectionString)
	localNode := Node.NewNode(p2pserver.configuration.BindingIPAddress, p2pserver.configuration.BindingPort)
	bytesRepresentation, err := json.Marshal(localNode)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(connectionString, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	//get nodes vom server
	var result []Node.Node
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println("get response from network connect request")
	log.Println(result)
	p2pserver.nodeBuffer.nodes = result

}

func (p2pserver *p2pserver) StartServer() {
	p2pserver.keepAlive()

	if strings.Compare(p2pserver.configuration.PeerServer, "") != 0 {
		log.Println("Connection String " + p2pserver.configuration.PeerServer + " found")
		//connection to server
		p2pserver.registerNetwork()
	} else {
		log.Println("Starting  new network")
	}

	log.Println("Server started on " + p2pserver.configuration.BindingIPAddress + ":" + p2pserver.configuration.BindingPort)

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/register", p2pserver.registerHandler)
	http.HandleFunc("/pushNode", p2pserver.pushNewNodeInfoHandler)
	http.HandleFunc("/config", p2pserver.configHandler)
	http.HandleFunc("/test", testHandler)

	http.ListenAndServe(":"+p2pserver.configuration.BindingPort, nil)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/ping")
	fmt.Fprintf(w, "OK")
}


func (p2pserver *p2pserver) configHandler(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Method)
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: %v", err)
		return
	}
	var result util.ConfigConnector
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	json.Unmarshal([]byte(buf.String()), &result)

	switch result.Command {
	case "LISTNODES":
		fmt.Fprintf(w, p2pserver.nodeBuffer.toJSON())
	default:
		fmt.Fprintf(w, "COMMAND INVALID")
	}
}

func (p2pserver *p2pserver) testHandler(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Method)

	fmt.Fprintf(w, p2pserver.payload.toJSON())
}


func (p2pserver *p2pserver) registerHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("register attempt")
	log.Println(r.Method)
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: %v", err)
		return
	}

	//parse client
	var resa Node.Node
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	json.Unmarshal([]byte(buf.String()), &resa)
	log.Println(resa)
	p2pserver.addNode(resa.IPaddress, resa.Port)

	//send to client
	log.Println("send: " + p2pserver.nodeBuffer.toJSON())
	fmt.Fprintf(w, p2pserver.nodeBuffer.toJSON())

	p2pserver.pushNode(resa.IPaddress, resa.Port)

}

func (p2pserver *p2pserver) pushNewNodeInfoHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("push new infos ...")
	log.Println(r.Method)
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: %v", err)
		return
	}

	//parse client
	var resa Node.Node
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	json.Unmarshal([]byte(buf.String()), &resa)
	log.Println(resa)
	p2pserver.addNode(resa.IPaddress, resa.Port)

	//send to client
	log.Println("send: " + "OK")
	fmt.Fprintf(w, "OK")

}

func (p2pserver *p2pserver) keepAlive() {
	ticker := time.NewTicker(time.Duration(p2pserver.configuration.KeepAlive) * 1000 * time.Millisecond)
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

					oldStamp := util.GetHourMinuteSecond(0, 0, -time.Duration(p2pserver.configuration.KeepAlive))
					if node.LastActive.Before(oldStamp) {

						log.Println("node too old")
						url := "http://" + node.IPaddress + ":" + node.Port + "/ping"
						requestData := util.GetRequest(url)
						if requestData == "OK" {
							log.Println("KeepAlive OK with " + url)
							node.SetActive()
						} else {
							log.Println("KeepAlive failed with " + url)
							p2pserver.nodeBuffer.deleteNode(node)
						}
					}
				}
			}
		}
	}()

}
