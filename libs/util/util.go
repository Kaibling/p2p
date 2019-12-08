package util

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/Kaibling/p2p/libs/Node"
)

type Configuration struct {
	BindingIPAddress string
	BindingPort      string
	PeerServer       string
	KeepAlive		 time.Duration
	NetworkName		 string
}

func ParseArguments() map[string]string {

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

/*
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
*/

func GetRequest(url string) string {

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

func ParseConfigurationFile(filepath string) *Configuration {

	if filepath == "config.json" {
		//default path found
		//create config file
		_, err := os.Stat(filepath)

		if os.IsNotExist(err) {
			log.Println("file does not exists")

			fo, err := os.Create("config.json")
			CheckError(err)

			returnConfig := new(Configuration)
			returnConfig.BindingIPAddress = "127.0.0.1"
			returnConfig.BindingPort = "7070"
			returnConfig.KeepAlive = 20
			returnConfig.NetworkName = "default"

			configString, err := json.Marshal(returnConfig)
			CheckError(err)

			_, err = fo.Write(configString)
			CheckError(err)

			defer fo.Close()
			log.Println("Configuration file created")
			return returnConfig
		}

	}
	log.Println("opening configuration file: " + filepath)
	returnConfig := new(Configuration)
	file, err := os.Open(filepath)
	CheckError(err)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&returnConfig)
	CheckError(err)
	return returnConfig
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func PostRequest(url string, postRequest []byte) string {

	log.Println("trying post request to: ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postRequest))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	resp.Body.Close()
	log.Println("response Body:", string(body))

	return string(body)
}

func FindNodeInArray(a []Node.Node, x Node.Node) int {
	for i, n := range a {
		if x.IPaddress == n.IPaddress && x.Port == n.Port {
			return i
		}
	}
	return len(a)
}

func GetHourMinuteSecond(hour, minute, second time.Duration) time.Time {
	return time.Now().Add(time.Hour*hour + time.Minute*minute + time.Second*second)
}
