package main

import "github.com/Kaibling/p2p/p2pServer"
import "github.com/Kaibling/p2p/libs/util"


import "encoding/json"
import "log"

type User struct {
Name string
LastUpdated string
}

	func (User *User) getVersion() string {
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



func main() {

	cliArguments := util.ParseArguments()
	config := util.ParseConfigurationFile(cliArguments["configFilePath"])
	server := p2pServer.Newp2pServer(config)

	userObject := new(User)
	server.AddPayload(userObject)
	server.StartServer()

}
