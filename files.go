package main

import (
	"fmt"
	"os"
	"strconv"
)

func commandFiles(args []string) {
	var action string
	if len(args) > 0 {
		action = args[0]
	} else {
		action = "help"
	}

	switch action {
	case "list":
		var initialFolder int64
		if len(args) > 0 {
			initialFolder, _ = strconv.ParseInt(args[0], 10, 64)
		}
		commandFilesList(initialFolder)

	case "link":
		var fileID int64
		if len(args) > 0 {
			fileID, _ = strconv.ParseInt(args[0], 10, 64)
		} else {
			fmt.Println("Error: please specify a file ID")
			os.Exit(1)
		}
		commandFilesGetLink(fileID)

	case "help":
		commandFilesHelp(0)

	default:
		fmt.Println("Error: unkonwn action '" + action + "' for command 'files'")
		commandFilesHelp(1)
	}
}

func commandFilesHelp(exitCode int) {
	fmt.Println(`Usage: pio [options] files <action>

Performs actions on files and directories in your Put.io storage.

Actions
  list [folder-id]   List files in specified folder. If no folder ID is given,
                     files are listed from the root folder.
  get-link           Display direct download link for given file.
  search             (not implemented)
  upload             (not implemented)
  new-folder         (not implemented)
  get                (not implemented)
  delete             (not implemented)
  rename             (not implemented)
  move               (not implemented)
  to-mp4             (not implemented)
  download           (not implemented)
  share              (not implemented)
  subtitles          (not implemented)
  playlist           (not implemented)

`)

	os.Exit(exitCode)
}

var treeDepth = 0

func commandFilesList(fileID int64) {
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
			commandFilesList(id)
			treeDepth--
		}
	}
}

func commandFilesGetLink(fileID int64) {
	link, err := putioAPI.GetDownloadLink(fileID)
	check(err)

	fmt.Println(link)
}
