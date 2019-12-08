package main

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "github.com/Kaibling/p2p/libs/util"
  "github.com/Kaibling/p2p/libs/Node"
  "encoding/json"
)

func main() {
startConsole()
//listNodes()
}


func startConsole() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("p2p Network")
	fmt.Println("------------")

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "q" || text == "quit" {
			return
		}
		switch text {
			case "q":
				return
			case "quit":
				return
			case "help":
				help()
			case "ls":
				listNodes()
			default:
				fmt.Println("unknown command")
		}

	}

}

func listNodes() {
	url := "http://localhost:7070/getNodes"
	result := util.GetRequest(url)
	if result == "NOK" {
		fmt.Println("no connection to server")
	} else {
		var nodes []Node.Node
		json.Unmarshal([]byte(result),&nodes)
		for _, node := range nodes {
			fmt.Printf("%s10 %s\n",node.IPaddress,node.Port)
		}
	}

}

func help() {
	fmt.Printf("ls     - list nodes\nq/quit - quit cli\nconfig - show network config\n")
}