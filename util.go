package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Configuration struct {
	BindingIPAddress string
	BindingPort      string
	PeerServer       string
}

func parseArguments() map[string]string {

	arguments := make(map[string]string)
	conString := flag.String("peerServer", "", "a connection string [ip/port]")
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

func getRequest(url string) string {

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return "NOK"
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return string(data)
}

func parseConfigurationFile(filepath string) *Configuration {

	if filepath == "config.json" {
		//default path found
		//create config file
		_, err := os.Stat(filepath)

		if os.IsNotExist(err) {
			log.Println("file does not exists")

			fo, err := os.Create("config.json")
			checkError(err)

			returnConfig := new(Configuration)
			returnConfig.BindingIPAddress = getPublicIP()
			returnConfig.BindingPort = "54321"
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

func postRequest(url string, postReqest []byte) {

	log.Println("trying post request to: ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postReqest))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
}

func Find(a []node, x node) int {
	for i, n := range a {
		if x.IPaddress == n.IPaddress && x.Port == n.Port {
			return i
		}
	}
	return len(a)
}
