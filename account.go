package main

import "fmt"

func commandAccount(args []string) {
	var action string
	if len(args) > 0 {
		action = args[0]
	} else {
		action = "info"
	}

	switch action {
	case "info":
		commandAccountInfo()

	default:
		fmt.Println("Error: unkonwn action '" + action + "' for command 'account'")
		commandAccountHelp(1)
	}
}

func commandAccountHelp(exitCode int) {
	fmt.Println("Todo: pio account help")
}

func commandAccountInfo() {
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
