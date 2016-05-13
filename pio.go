package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"

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

func send(t string, uri string, arguments string) (*jason.Object, error) {
	url := "https://api.put.io/v2" + uri + "?oauth_token=" + t + "&" + arguments
	response, err := http.Get(url)
	check(err)
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	check(err)

	v, err := jason.NewObjectFromBytes([]byte(contents))
	return v, err
}

var treeDepth = 0

func list(t string, fileID int64) {
	v, err := send(t, "/files/list", "parent_id="+strconv.FormatInt(fileID, 10))
	check(err)

	status, _ := v.GetString("status")
	if status == "ERROR" {
		fmt.Println(`
Error during put.io API access, please check your internet connectivity and your
token configuration. See --help for more info.`)
		os.Exit(1)
	}

	// Tree formatting
	spaces := ""
	for i := 0; i < treeDepth; i++ {
		spaces = spaces + "\t"
	}

	files, _ := v.GetObjectArray("files")
	for _, file := range files {
		name, _ := file.GetString("name")
		contentType, _ := file.GetString("content_type")
		id, _ := file.GetInt64("id")
		fmt.Printf("%s%d: [%s] %s\n", spaces, id, contentType, name)

		if contentType == "application/x-directory" {
			treeDepth++
			list(t, id)
			treeDepth--
		}
	}
}

func info(t string) {
	fmt.Println("Account information ...")

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
		var initialFolder int64
		if len(os.Args) > 2 {
			initialFolder, _ = strconv.ParseInt(os.Args[2], 10, 64)
		}
		list(config.Auth.Token, initialFolder)
	case "info":
		info(config.Auth.Token)
	}

}
