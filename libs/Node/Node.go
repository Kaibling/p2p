package Node

import (
	"encoding/json"
	"log"
	"time"
)

type Node struct {
	IPaddress  string
	Port       string
	LastActive time.Time
}

func NewNode(ip string, port string) Node {

	returnNode := new(Node)
	returnNode.IPaddress = ip
	returnNode.Port = port
	returnNode.LastActive = time.Now()
	return *returnNode
}
func (node Node) toJSONString() string {
	nodeString, err := json.Marshal(node)
	if err != nil {
		log.Fatalln(err)
	}
	return string(nodeString)
}

func (node *Node) SetActive() {
	node.LastActive = time.Now()
}
