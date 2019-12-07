package main


func main() {

    
	//testa()
	//go startConsole()

    cliArguments := parseArguments()
	config := parseConfigurationFile(cliArguments["configFilePath"])
	server := newP2Pserver(config)
	server.startServer()

}