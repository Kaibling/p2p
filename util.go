package main

import ("flag"
"net/http"
"io/ioutil"
"log"
"os"
"encoding/json"
)

type Configuration struct {
	IPAddress string
	Port string
}

func parseArguments() map[string] string{

	arguments := make( map[string] string)

	conString := flag.String("server", "", "a connection string [ip/port]")
	configString := flag.String("config", "config.json", "filepath to configuration file")
	flag.Parse()

    arguments["connectionstring"] = *conString
    arguments["configFilePath"] = *configString

	log.Print("load command line arguments ")
	log.Print(arguments)
	return arguments
}

func getPublicIP() string {
    url := "https://api.ipify.org?format=text"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
    }
    return string(ip)
}

func parseConfigurationFile(filepath string) *Configuration{

    if filepath == "config.json" {
        //default path found
        //create config file

        _, err := os.Stat(filepath)
        if os.IsNotExist(err) {
			log.Println("file does not exists")
			
			fo, err := os.Create("config.json")
			checkError(err)
			
			returnConfig := new(Configuration)
			returnConfig.IPAddress = getPublicIP()
			returnConfig.Port = "54321"
			configString, err := json.Marshal(returnConfig)
			checkError(err)
			_, err = fo.Write(configString)
			checkError(err)

			defer fo.Close()
			log.Println("Configuration file created")
			return returnConfig
        }
       
    }
	log.Println("opening configuration file: " + filepath)
	returnConfig := new(Configuration)
	file, err := os.Open(filepath)
	checkError(err)
	decoder := json.NewDecoder(file) 
	err = decoder.Decode(&returnConfig)
	checkError(err)
	return returnConfig
}

    func checkError(err error) {
        if err != nil {
			log.Fatalln(err)
		}
    }