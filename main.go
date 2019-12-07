package main

func main() {

        cliArguments := parseArguments()
        config := parseConfigurationFile(cliArguments["configFilePath"])
        server := newP2Pserver(config)
        server.startServer()

}