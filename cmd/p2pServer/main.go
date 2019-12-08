package main

import "github.com/Kaibling/p2p/p2pServer"
import "github.com/Kaibling/p2p/libs/util"


func main() {

	cliArguments := util.ParseArguments()
	config := util.ParseConfigurationFile(cliArguments["configFilePath"])
	server := p2pServer.Newp2pServer(config)
	server.StartServer()

}
