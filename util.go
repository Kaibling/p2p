package main

import ("flag"
"net/http"
"io/ioutil"
)

func parseArguments() string{
	connectionstring := *flag.String("server", "", "a connection string")
	flag.Parse()
	return connectionstring

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