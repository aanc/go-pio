package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
)

func commandTransfers(args []string) {
	var action string
	if len(args) > 0 {
		action = args[0]
	} else {
		action = listAction
	}

	switch action {
	case listAction:
		commandTransfersList()

	case helpAction:
		commandTransfersHelp(0)

	default:
		fmt.Println("Error: unkonwn action '" + action + "' for command 'transfers'")
		commandTransfersHelp(1)
	}
}

func commandTransfersHelp(exitCode int) {
	os.Exit(exitCode)
}

func commandTransfersList() {
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
