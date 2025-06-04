package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	records, err := New()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer records.db.Close()

	fmt.Println("Welcome to the password manager! 🔑")

	for {
		menuView()
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			if err := saveNewRecordView(records); err != nil {
				errorView(err)
			}
		case 2:
			if err := getRecordView(records); err != nil {
				errorView(err)
			}
		case 3:
			if err := updateRecordView(records); err != nil {
				errorView(err)
			}
		case 4:
			if err := deleteRecordView(records); err != nil {
				errorView(err)
			}
		case 5:
			if err := listDomainsView(records); err != nil {
				errorView(err)
			}
		case 6:
			exitView()
		default:
			errorView(fmt.Errorf("invalid choice. try again"))
		}
	}
}
