package main
import ("net/http"
"fmt"
"log"
"encoding/json"
"strings"
"bytes"
)

type p2pserver struct {
    nodes []node
	publicIP string
	configuration *Configuration
}

func newP2Pserver(configuration *Configuration) *p2pserver {

	returnP2Pserver := new(p2pserver)
	returnP2Pserver.publicIP = getPublicIP()
	returnP2Pserver.configuration = configuration
    return returnP2Pserver
    
}

func (p2pserver *p2pserver) addNode(ipAddress string, port string) {

    newNode := newNode(ipAddress,port)
    p2pserver.nodes = append(p2pserver.nodes, newNode)
    log.Print("add Node to Buffer: ")
    log.Println(p2pserver.nodes)
    
}

func (p2pserver *p2pserver) registerNetwork() {

    connectionString := "http://" + p2pserver.configuration.PeerServer+"/register"
    log.Println("trying to register to " + connectionString)
	localNode := newNode(p2pserver.configuration.BindingIPAddress,p2pserver.configuration.BindingPort)
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
    p2pserver.nodes = result
    


}

func (p2pserver *p2pserver) startServer() {

    if strings.Compare(p2pserver.configuration.PeerServer, "") != 0 {
        log.Println("Connection String " + p2pserver.configuration.PeerServer + " found")
		//connection to server
		p2pserver.registerNetwork()
    } else {
        log.Println("Starting  new network")
        p2pserver.addNode(p2pserver.configuration.BindingIPAddress,p2pserver.configuration.BindingPort)
    }
	
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {

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
        p2pserver.addNode(resa.Ipaddress,resa.Port)
        //send to client
        jnodes,err := json.Marshal(p2pserver.nodes)
		if err != nil {
			log.Fatalln(err)
        }
		log.Println("send: " + string(jnodes))
        fmt.Fprintf(w, string(jnodes))
		
	})
	http.HandleFunc("/getNodes", func(w http.ResponseWriter, r *http.Request) {
        
        log.Println(r)
		jnodes,err := json.Marshal(p2pserver.nodes)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Fprintf(w, string(jnodes))
	})
	log.Println("Server started on " + p2pserver.configuration.BindingIPAddress + ":" + p2pserver.configuration.BindingPort)
	http.ListenAndServe(":"+p2pserver.configuration.BindingPort, nil)
}
