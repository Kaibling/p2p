package main
import ("net/http"
"fmt"
"log"
"encoding/json"
"strings"
"bytes"
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

func (p2pserver p2pserver) registerNetwork(connectionString string) {

	localNode := newNode("1.2.3.4","1243")
	bytesRepresentation, err := json.Marshal(localNode)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(connectionString, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
	//log.Println(result["data"])
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
		log.Println("send: " + string(jnodes))
		for _,val := range p2pserver.nodes {
			log.Println(val.toJSONString())

		}
		fmt.Fprintf(w, string(jnodes))

		switch r.Method {
    case "GET":     
         fmt.Println("get")
	case "POST":
		fmt.Println("form")
		}

// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
        if err := r.ParseForm(); err != nil {
            fmt.Fprintf(w, "ParseForm() err: %v", err)
            return
        }
        fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		/*
		var resa node
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		log.Println("recive string " + buf.String())

		log.Println("-----node")
		json.Unmarshal([]byte(buf.String()), &resa)
		if err != nil {
			log.Println("node error")
			log.Println(err)
			log.Printf("%+v\n", resa)
		}

		log.Println("-----")
		log.Println(resa)
		log.Println("-----")
		*/
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
