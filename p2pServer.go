package main
import ("net/http"
"fmt"
"log"
"encoding/json"
"strings"
)

type p2pserver struct {
    nodes map[string]node
    bindingPort string
    bindingIP string
    publicIP string
}

func newP2Pserver(bindingIP string , bindingPort string) *p2pserver {

	returnP2Pserver := new(p2pserver)
    returnP2Pserver.nodes = make(map[string]node)
    returnP2Pserver.bindingPort = bindingPort
    returnP2Pserver.bindingIP = bindingIP
    returnP2Pserver.publicIP = getPublicIP()
	return returnP2Pserver
}

func (p2pserver p2pserver) addNode(ipAddress string, port string) {
	newNode := newNode(ipAddress,port)
	p2pserver.nodes[newNode.Ipaddress] = *newNode
}

func (p2pserver p2pserver) startServer(connectionString string) {

    if strings.Compare(connectionString, "") != 0 {
        log.Println("Connection String " + connectionString + " found")
        log.Println("Trying to connect")
        //connection to server
    } else {
        log.Println("Starting  new network")
        p2pserver.addNode(p2pserver.publicIP,p2pserver.bindingPort)
    }
	
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {

		jnodes,err := json.Marshal(p2pserver.nodes)
		if err != nil {
			log.Fatalln(err)
        }
        fmt.Println(string(jnodes))
        fmt.Fprintf(w, string(jnodes))
        var result map[string]interface{}
        json.Unmarshal([]byte(r.Body), &result)
        fmt.Println(result)
	})
	http.HandleFunc("/getNodes", func(w http.ResponseWriter, r *http.Request) {
		
		jnodes,err := json.Marshal(p2pserver.nodes)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Fprintf(w, string(jnodes))
	})
	log.Println("Server started")
	http.ListenAndServe(":"+p2pserver.bindingPort, nil)
}
