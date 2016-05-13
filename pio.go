package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/antonholmquist/jason"
	"gopkg.in/gcfg.v1"
)

// Configuration : used to store config retrieved from config file
type Configuration struct {
	Auth struct {
		Token string `gcfg:"token"`
	}
}

func check(e error) {
	if e != nil {
		fmt.Printf("%s", e)
		panic(e)
	}
}

func initAndReadConfigFile(f string) string {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		fmt.Println("Initializing config file in " + f)
		cfgInit := `[auth]
token=
`
		err := ioutil.WriteFile(f, []byte(cfgInit), 0600)
		check(err)
	}

	configFileContent, err := ioutil.ReadFile(f)
	check(err)
	return string(configFileContent)
}

func main() {
	// Get home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	config := Configuration{}

	err = gcfg.ReadStringInto(&config, initAndReadConfigFile(usr.HomeDir+"/.piorc"))
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err)
	}

	if config.Auth.Token == "" {
		fmt.Println("Please edit the file " + usr.HomeDir + "/.piorc with your OAuth token.")
		fmt.Println("To get this token, go to http://aanc.github.io/go-pio")
		os.Exit(1)
	}

	// Action routing
	action := os.Args[1]
	switch action {
	case "list":
		list(config.Auth.Token)
	case "info":
		info(config.Auth.Token)
	}

}

func list(t string) {
	response, err := http.Get("https://api.put.io/v2/files/list?oauth_token=" + t)
	check(err)
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	check(err)

	v, err := jason.NewObjectFromBytes([]byte(contents))
	check(err)

	status, _ := v.GetString("status")
	if status == "ERROR" {
		fmt.Println("Error during put.io API access, please check your internet connectivity and your token configuration.")
		fmt.Println("See --help for more info")
		os.Exit(1)
	}

	files, _ := v.GetObjectArray("files")
	for i, file := range files {
		name, _ := file.GetString("name")
		contentType, _ := file.GetString("content_type")
		fmt.Printf("%.2d: [%s] %s\n", i, contentType, name)
	}
}

func info(t string) {
	fmt.Println("Account information ...")

}
