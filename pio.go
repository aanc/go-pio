package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

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
	fmt.Printf("Using token %s\n", config.Auth.Token)

}
