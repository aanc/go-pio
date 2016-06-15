package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
)

func commandEvents(args []string) {
	var action string
	if len(args) > 0 {
		action = args[0]
	} else {
		action = listAction
	}

	switch action {
	case listAction:
		commandEventsList()

	case helpAction:
		commandEventsHelp(0)

	default:
		fmt.Println("Error: unkonwn action '" + action + "' for command 'transfers'")
		commandEventsHelp(1)
	}
}

func commandEventsHelp(exitCode int) {
	os.Exit(exitCode)
}

func commandEventsList() {
	json, err := putioAPI.Events()
	check(err)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)

	fmt.Fprint(w, "Date")
	fmt.Fprint(w, "\tEvent")
	fmt.Fprint(w, "\tDetails")
	fmt.Fprintln(w, "")

	events, _ := json.GetObjectArray("events")

	for _, event := range events {
		eventType, _ := event.GetString("type")
		eventDate, _ := event.GetString("created_at")

		fmt.Fprintf(w, "%s", eventDate)

		switch eventType {
		case "transfer_completed":
			// keys: created_at, file_id, id, transfer_name, transfer_size, type

			transferName, _ := event.GetString("transfer_name")

			fmt.Fprint(w, "\tCOMPLETED")
			fmt.Fprintf(w, "\t%s", transferName)

		case "zip_created":
			// Keys: type, zip_id, zip_size, created_at, id

			zipID, _ := event.GetInt64("zip_id")
			zipSize, _ := event.GetInt64("zip_size")

			fmt.Fprint(w, "\tZIPPED    ")
			fmt.Fprintf(w, "\tID: %d (%s)", zipID, humanize.Bytes(uint64(zipSize)))

		default:
			fmt.Println("Ahem... It seems you encountered an unsupported event type, sorry !")
			fmt.Println("Please open an issue at https://github.com/aanc/go-pio/issues with the following output:")

			fmt.Println("Event type: " + eventType)
			for key := range event.Map() {
				fmt.Println("Key: " + key)
			}
			fmt.Println("")

		}
		fmt.Fprintln(w, "")
		w.Flush()
	}

}
