package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"

	"github.com/aanc/pio/putio"

	"gopkg.in/gcfg.v1"
)

var putioAPI putio.Config

// Configuration : used to store config retrieved from config file
type Configuration struct {
	Auth struct {
		Token string `gcfg:"token"`
	}
}

func usage(exitCode int) {
	fmt.Println(`Usage: pio [options] <command> args ...

Options:
  ...

Commands:
  list      List files
  config    Update configuration file
  search    Search for files matching the given word

Run 'pio <command> --help' for more information about a command.`)

	os.Exit(exitCode)
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

var treeDepth = 0

func commandList(fileID int64) {
	v, err := putioAPI.List(fileID)
	check(err)

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
			commandList(id)
			treeDepth--
		}
	}
}

func commandInfo() {
	info, err := putioAPI.AccountInfo()
	check(err)

	username, _ := info.GetString("info", "username")
	fmt.Println("  Username: " + username)
}

func main() {
	// Get home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Read configuration from file
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

	// Initialize put.io API
	putioAPI.SetToken(config.Auth.Token)

	// Action routing
	if len(os.Args) <= 1 {
		usage(0)
	}

	action := os.Args[1]
	switch action {
	case "list":
		var initialFolder int64
		if len(os.Args) > 2 {
			initialFolder, _ = strconv.ParseInt(os.Args[2], 10, 64)
		}
		commandList(initialFolder)

	case "info":
		commandInfo()
	}
}
