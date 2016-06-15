package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/aanc/pio/putio"

	"gopkg.in/gcfg.v1"
)

const listAction = "list"
const helpAction = "help"

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
  --token=<token>     Use that token for authentication with Put.io's API,
                      instead of the one configured in ~/.piorc (not
                      implemented)

  --conf-file=<file>  Use the specified file for configuration, instead of the
                      default ~/.piorc (not implemented)

Commands:

  files                Manage files and folders
      list
      search           (not implemented)
      upload           (not implemented)
      new-folder       (not implemented)
      get              (not implemented)
      delete           (not implemented)
      rename           (not implemented)
      move             (not implemented)
      to-mp4           (not implemented)
      download         (not implemented)
      share            (not implemented)
      subtitles        (not implemented)
      playlist         (not implemented)

  events               Manage dashboard events
      list
      clear            (not implemented)

  transfers            Manage active or finished transfers
      list             (not implemented)
      add              (not implemented)
      get              (not implemented)
      retry            (not implemented)
      cancel           (not implemented)
      clean            (not implemented)

  friends              Manage friends
      list             (not implemented)
      requests         (not implemented)
      send-request     (not implemented)
      approve          (not implemented)
      deny             (not implemented)
      unfriend         (not implemented)

  account              Manage account
      info
      settings         (not implemented)

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

	command := os.Args[1]
	otherArgs := []string{}
	if len(os.Args) > 1 {
		otherArgs = os.Args[2:]
	}

	switch command {
	case "files":
		commandFiles(otherArgs)

	case "account":
		commandAccount(otherArgs)

	case "transfers":
		commandTransfers(otherArgs)

	case "events":
		commandEvents(otherArgs)

	default:
		fmt.Println("Error: unkonwn command '" + command + "'")
		usage(2)
	}
}
