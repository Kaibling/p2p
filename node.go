package main


type node struct {
	Ipaddress  string
	Port       string
	LastActive int64
}

func newNode(ip string, port string) *node {

	returnNode := new(node)
	returnNode.Ipaddress = ip
	returnNode.Port = port
	return returnNode
}