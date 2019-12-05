package main


func main() {

    //go startConsole()
    connectionString := parseArguments()

	server := newP2Pserver("123.123.123.1","5421")
	server.startServer(connectionString)
}