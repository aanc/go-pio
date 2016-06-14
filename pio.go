package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"text/tabwriter"

	"github.com/aanc/pio/putio"
	"github.com/dustin/go-humanize"

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
  config       Update configuration file
  link         Get file download link
  info         Display information about configured put.io account
  list         List files
  search       Search for files matching the given word
  transfers    Get transferts list

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
	mail, _ := info.GetString("info", "mail")
	planExpirationDate, _ := info.GetString("info", "plan_expiration_date")
	diskAvail, _ := info.GetInt64("info", "disk", "avail")
	diskUsed, _ := info.GetInt64("info", "disk", "used")
	diskSize, _ := info.GetInt64("info", "disk", "size")

	diskAvail = diskAvail / 1024 / 1024 / 1024
	diskUsed = diskUsed / 1024 / 1024 / 1024
	diskSize = diskSize / 1024 / 1024 / 1024

	fmt.Printf(`Account information:

  Username: %s
  Mail: %s
  Disk: %.2d/%.2dGB (%.2dGB free)
  Expiration: %s
`, username, mail, diskUsed, diskSize, diskAvail, planExpirationDate)
}

func commandTransfers() {
	json, err := putioAPI.Transfers()
	check(err)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)

	fmt.Fprint(w, "ID")
	fmt.Fprint(w, "\tFile ID")
	fmt.Fprint(w, "\tFile name")
	fmt.Fprint(w, "\tStatus")
	fmt.Fprint(w, "\tSize")
	fmt.Fprint(w, "\tDown")
	fmt.Fprint(w, "\tUp")
	fmt.Fprint(w, "\tRatio")
	fmt.Fprintln(w, "")

	transfers, _ := json.GetObjectArray("transfers")
	for _, t := range transfers {
		name, _ := t.GetString("name")
		status, _ := t.GetString("status")
		id, _ := t.GetInt64("id")
		downSpeed, _ := t.GetInt64("down_speed")
		upSpeed, _ := t.GetInt64("up_speed")
		size, _ := t.GetInt64("size")
		ratio, _ := t.GetFloat64("current_ratio")
		fileID, _ := t.GetInt64("file_id")

		fmt.Fprint(w, strconv.FormatInt(id, 10))
		fmt.Fprint(w, "\t"+strconv.FormatInt(fileID, 10))
		fmt.Fprint(w, "\t"+name)
		fmt.Fprint(w, "\t"+status)
		fmt.Fprint(w, "\t"+humanize.Bytes(uint64(size)))
		fmt.Fprint(w, "\t"+humanize.Bytes(uint64(downSpeed))+"/s")
		fmt.Fprint(w, "\t"+humanize.Bytes(uint64(upSpeed))+"/s")
		fmt.Fprint(w, "\t"+strconv.FormatFloat(ratio, 'f', -1, 64))
		fmt.Fprintln(w, "")

	}
	w.Flush()
}

func commandDlLink(fileID int64) {
	link, err := putioAPI.GetDownloadLink(fileID)
	check(err)

	fmt.Println(link)
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
	switch command {
	case "list":
		var initialFolder int64
		if len(os.Args) > 2 {
			initialFolder, _ = strconv.ParseInt(os.Args[2], 10, 64)
		}
		commandList(initialFolder)

	case "info":
		commandInfo()

	case "transfers":
		commandTransfers()

	case "link":
		var fileToDL int64
		if len(os.Args) > 2 {
			fileToDL, _ = strconv.ParseInt(os.Args[2], 10, 64)
		}
		commandDlLink(fileToDL)

	default:
		fmt.Println("Error: unkonwn command '" + command + "'")
		usage(2)
	}

}
