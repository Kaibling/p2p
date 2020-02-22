package peerserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Kaibling/p2p/libs/Node"
	"github.com/Kaibling/p2p/libs/util"
	"log"
	"net/http"
	"strings"
	"time"
)

func init(){
    //log.SetPrefix("TRACE: ")
    log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
    log.Println("Logging started")
}

type payload struct {
	Version string
	Data    interface{}
}

//Peerserver da
type Peerserver struct {
	nodes         []Node.Node
	publicIP      string
	configuration *util.Configuration
	payload       *payload
	A             string
}

func (p2pserver *Peerserver) addNode(ipAddress string, port string) {
	node := Node.NewNode(ipAddress, port)
	i := util.FindNodeInArray(p2pserver.nodes, node)
	if i == len(p2pserver.nodes) {
		//.Println("new element found")
		p2pserver.nodes = append(p2pserver.nodes, node)
		//log.Print("add Node to Buffer: ")
		//log.Println(node)

	} else {
		//log.Println("node already in buffer. skip")
	}

}

func (p2pserver *Peerserver) deleteNode(ipAddress string, port string) {
	node := Node.NewNode(ipAddress, port)
	i := util.FindNodeInArray(p2pserver.nodes, node)
	if i == len(p2pserver.nodes) {
		//log.Println("element not found")
		return
	}
	p2pserver.nodes = append(p2pserver.nodes[:i], p2pserver.nodes[i+1:]...)
	log.Println(node)
	//log.Println("node deleted")

}

func (p2pserver *Peerserver) nodesToJSON() string {
	jnodes, err := json.Marshal(p2pserver.nodes)
	if err != nil {
		log.Fatalln(err)
	}
	return string(jnodes)
}

//Newpeerserver constructor
func Newpeerserver(configuration *util.Configuration) *Peerserver {

	returnP2Pserver := new(Peerserver)
	returnP2Pserver.publicIP = "undef"
	returnP2Pserver.configuration = configuration
	returnP2Pserver.addNode(configuration.BindingIPAddress, configuration.BindingPort)
	return returnP2Pserver

}

func (p2pserver *Peerserver) pushNode(ipAddress string, port string) {

	newNode := Node.NewNode(ipAddress, port)
	//push to network
	for _, node := range p2pserver.nodes {
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

//GetPayloadData recieve Data
func (p2pserver *Peerserver) GetPayloadData() interface{} {
	return p2pserver.payload.Data
}

//GeneratePayload add an object
func (p2pserver *Peerserver) GeneratePayload(data interface{}) {
	p2pserver.payload = &payload{
		Version: "0",
		Data:    data,
	}
	bytesRepresentation, err := json.Marshal(p2pserver.payload)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("GeneratePayload: payload saved: ", string(bytesRepresentation))
}

func (p2pserver *Peerserver) registerNetwork() {

	//send own node data to server
	connectionString := "http://" + p2pserver.configuration.PeerServer + "/register"
	log.Println("registerNetwork: trying to register to " + connectionString)
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
	log.Println("registerNetwork: getting network topology",result)
	p2pserver.nodes = result

	//get Payload
	url := "http://" + p2pserver.configuration.PeerServer + "/receivePL"
	requestData := util.GetRequest(url)
	var payloadresult payload
	json.Unmarshal([]byte(requestData), &payloadresult)
	log.Println("registerNetwork: payload recieved from  from " ,url, " -> ",payloadresult.Version)
	*p2pserver.payload = payloadresult

}

//StartServer start server
func (p2pserver *Peerserver) StartServer() {

	if strings.Compare(p2pserver.configuration.PeerServer, "") != 0 {
		log.Println("StartServer: Connection String " + p2pserver.configuration.PeerServer + " found")
		//connection to server
		p2pserver.registerNetwork()
	} else {
		log.Println("StartServer: Starting  new network")
	}

	log.Println("StartServer: Server started on " + p2pserver.configuration.BindingIPAddress + ":" + p2pserver.configuration.BindingPort)

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/register", p2pserver.registerHandler)
	http.HandleFunc("/pushNode", p2pserver.pushNewNodeInfoHandler)
	http.HandleFunc("/config", p2pserver.configHandler)
	http.HandleFunc("/receivePL", p2pserver.receivePayloadHandler)
	//http.HandleFunc("/pullPL", p2pserver.pullPayloadHandler)
	http.HandleFunc("/health", p2pserver.healthHandler)

	http.ListenAndServe(":"+p2pserver.configuration.BindingPort, nil)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("pingHandler: /ping")
	fmt.Fprintf(w, "OK")
}

func (p2pserver *Peerserver) healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("healthHandler: health check")
	fmt.Fprintf(w, "connected Nodes: %s\n", p2pserver.nodesToJSON())
	bytesRepresentation, err := json.Marshal(p2pserver.payload)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(w, "payload: %s\n", string(bytesRepresentation))
}

func (p2pserver *Peerserver) configHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("configHandler:",r.Method)
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
		fmt.Fprintf(w, p2pserver.nodesToJSON())
	default:
		fmt.Fprintf(w, "COMMAND INVALID")
	}
}

//func (p2pserver *Peerserver) pullPayloadHandler(w http.ResponseWriter, r *http.Request) {
//
//	log.Println(r.Method)
//	//parse client
//	var resa payload
//	buf := new(bytes.Buffer)
//	buf.ReadFrom(r.Body)
//	json.Unmarshal([]byte(buf.String()), &resa)
//	log.Println(resa)
//	p2pserver.payload = &resa
//}

func (p2pserver *Peerserver) receivePayloadHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("receivePayloadHandler",r.Method)
	//parse client
	bytesRepresentation, err := json.Marshal(p2pserver.payload)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(w, string(bytesRepresentation))
}

func (p2pserver *Peerserver) registerHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("registerHandler: register attempt")
	log.Println("registerHandler:",r.Method)
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
	log.Println("registerHandler: send: " + p2pserver.nodesToJSON())
	fmt.Fprintf(w, p2pserver.nodesToJSON())

	p2pserver.pushNode(resa.IPaddress, resa.Port)

}

func (p2pserver *Peerserver) pushNewNodeInfoHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("pushNewNodeInfoHandler: push new infos .... Request with ",r.Method)
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
	log.Println("pushNewNodeInfoHandler: send: " + "OK")
	fmt.Fprintf(w, "OK")

}

func (p2pserver *Peerserver) keepAlive() {
	ticker := time.NewTicker(time.Duration(p2pserver.configuration.KeepAlive) * 1000 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for _, node := range p2pserver.nodes {
					if node.IPaddress == p2pserver.configuration.BindingIPAddress && node.Port == p2pserver.configuration.BindingPort {
						continue
					}

					oldStamp := util.GetHourMinuteSecond(0, 0, -time.Duration(p2pserver.configuration.KeepAlive))
					if node.LastActive.Before(oldStamp) {

						log.Println("keepAlive: node too old")
						url := "http://" + node.IPaddress + ":" + node.Port + "/ping"
						requestData := util.GetRequest(url)
						if requestData == "OK" {
							log.Println("keepAlive: KeepAlive OK with " + url)
							node.SetActive()
						} else {
							log.Println("keepAlive: KeepAlive failed with " + url)
							p2pserver.deleteNode(node.IPaddress, node.Port)
						}
					}
				}
			}
		}
	}()

}
