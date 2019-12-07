package main

import ("encoding/json"
		"log"			
)

type node struct {
	Ipaddress  string
	Port       string
	LastActive int64
}

func newNode(ip string, port string) node {

	returnNode := new(node)
	returnNode.Ipaddress = ip
	returnNode.Port = port
	return *returnNode
}
func (node node) toJSONString() string {
	nodeString,err := json.Marshal(node)
		if err != nil {
			log.Fatalln(err)
        }
    return string(nodeString)
}