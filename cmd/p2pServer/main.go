package main

import "github.com/Kaibling/p2p/peerserver"
import "github.com/Kaibling/p2p/libs/util"


import "encoding/json"
import "log"

type User struct {
Name string
LastUpdated int
}

	func (User *User) getVersion() int {
		return User.LastUpdated

	}
	func (User *User) saveData(name interface {}) {
		User.Name = name.(string)

	}
	func (User *User) getData() interface {} {
		return User.Name

	}
	func (User *User) toJSON() string {
		jsonByte, err := json.Marshal(User.Name)
		if err != nil {
			log.Println(err)
		}
		return string(jsonByte)
	}

	
func init(){
    //log.SetPrefix("TRACE: ")
    log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
    log.Println("Logging started")
}


func main() {

	

	cliArguments := util.ParseArguments()
	config := util.ParseConfigurationFile(cliArguments["configFilePath"])
	server := peerserver.Newpeerserver(config)

	userObject := &User{
        Name: "Hans",
        LastUpdated: 123,
    }
	server.GeneratePayload(userObject)
	server.StartServer()

}
