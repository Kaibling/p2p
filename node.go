package main

import ("encoding/json"
        "log"	
        "time"		
)

type node struct {
	IPaddress  string
	Port       string
	LastActive time.Time
}

func newNode(ip string, port string) node {

	returnNode := new(node)
	returnNode.IPaddress = ip
    returnNode.Port = port
    returnNode.LastActive =  time.Now()
	return *returnNode
}
func (node node) toJSONString() string {
	nodeString,err := json.Marshal(node)
		if err != nil {
			log.Fatalln(err)
        }
    return string(nodeString)
}

func (node *node) setActive() {
    node.LastActive = time.Now()
}