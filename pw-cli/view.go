package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

var (
	ASCII = `
		
	██████╗ ██╗    ██╗ ██████╗ ██╗     ██╗     
	██╔══██╗██║    ██║██╔════╝ ██║     ██║     
	██████╔╝██║ █╗ ██║██║      ██║     ██║     
	██╔═══╝ ██║███╗██║██║      ██║     ██║     
	██║     ╚███╔███╔╝╚██████╗ ███████╗██║ 
	╚═╝      ╚══╝╚══╝  ╚═════╝ ╚══════╝╚═╝

         PWCli - Secure & Simple 🔐

	`
	menu = fmt.Sprintf(`
	1. Save a new record
	2. Get a record
	3. Update a record
	4. Delete a record
	5. List all domains
	6. Exit
	`+"\n%s\nEnter your choice: ", sectionDivisor)
	sectionDivisor = "----------------------------------------"
)

func typePassword() (string, error) {
	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	return string(passwordBytes), nil

}

func randomGeneratePassword(length int) (string, error) {
	if length < 4 {
		return "", fmt.Errorf("password length must be at least 4")
	}

	lowercase := "abcdefghijklmnopqrstuvwxyz"
	uppercase := strings.ToUpper(lowercase)
	digits := "0123456789"
	symbols := `"'#%&*()_-=+{}[]^~;:.><,/!\/|`
	charset := lowercase + uppercase + digits + symbols

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		var password []byte
		for len(password) < length {
			password = append(password, charset[r.Intn(len(charset))])
		}

		passStr := string(password)
		if strings.ContainsAny(passStr, lowercase) &&
			strings.ContainsAny(passStr, uppercase) &&
			strings.ContainsAny(passStr, digits) &&
			strings.ContainsAny(passStr, symbols) &&
			!charRepeatOverNTimes(passStr, 2) {
			return passStr, nil
		}
	}
}

func charRepeatOverNTimes(str string, times int) bool {
	counter := make(map[rune]int)
	for _, l := range str {
		counter[rune(l)]++
		if counter[rune(l)] > times {
			return true
		}
	}
	return false
}

func generatePassword() (string, error) {
	var choice string
	fmt.Print(`
		Choose how you want to generate password:
			1. Random Generate Password
			2. Type Password
		`)
	fmt.Scanln(&choice)
	switch choice {
	case "1":
		return randomGeneratePassword(10)
	case "2":
		return typePassword()
	default:
		return "", fmt.Errorf("invalid option, try again")
	}

}

func saveNewRecordView(records Storage) error {
	var domain, user, password string
	fmt.Print("Enter domain: ")
	fmt.Scanln(&domain)
	fmt.Print("Enter user: ")
	fmt.Scanln(&user)
	password, err := generatePassword()
	if err != nil {
		return fmt.Errorf("generating password error. try again %s", err)
	}
	record := NewRecord(user, password)
	if err := records.SaveRecord(domain, record); err != nil {
		return fmt.Errorf("saving record error, try again %s", err)
	}
	fmt.Println("Record saved successfully")
	fmt.Println(sectionDivisor)
	fmt.Println("Press Enter to continue...")
	fmt.Scanln()
	return nil
}

func getRecordView(records Storage) error {
	var domain string
	fmt.Print("Enter domain: ")
	fmt.Scanln(&domain)
	record, err := records.GetRecord(domain)
	if err != nil {
		return fmt.Errorf("Record not found... Try again... %s", err)
	}
	fmt.Println("User:", record.User)
	fmt.Println("Password:", record.Password)
	fmt.Println(sectionDivisor)
	fmt.Println("Press Enter to continue...")
	fmt.Scanln()
	return nil
}

func updateRecordView(records Storage) error {
	var domain, user, password string
	fmt.Print("Enter domain: ")
	fmt.Scanln(&domain)
	fmt.Print("Enter user: ")
	fmt.Scanln(&user)
	password, err := generatePassword()
	if err != nil {
		return fmt.Errorf("generating password error. Try again %s", err)
	}
	record := NewRecord(user, password)
	if err = records.UpdateRecord(domain, record); err != nil {
		return fmt.Errorf("saving record error, try again %s", err)
	}
	fmt.Println("Record updated successfully")
	fmt.Println(sectionDivisor)
	fmt.Println("Press Enter to continue...")
	fmt.Scanln()
	return nil
}

func deleteRecordView(records Storage) error {
	var domain string
	fmt.Print("Enter domain: ")
	fmt.Scanln(&domain)
	if err := records.DeleteRecord(domain); err != nil {
		fmt.Errorf("deleting record error. Try again %s", err)
	}
	fmt.Println("Record deleted successfully")
	fmt.Println(sectionDivisor)
	fmt.Println("Press Enter to continue...")
	fmt.Scanln()
	return nil
}

func listDomainsView(records Storage) error {
	domains, err := records.ListDomains()
	if err != nil {
		fmt.Errorf("listing domains error. try again %s", err)
	}
	fmt.Println("Domains:")
	fmt.Println("")
	for i, domain := range domains {
		fmt.Printf("%d: %s\n", i+1, domain)
	}
	fmt.Println(sectionDivisor)
	fmt.Println("Press Enter to continue...")
	fmt.Scanln()
	return nil
}

func errorView(err error) {
	fmt.Println(err)
	fmt.Println(sectionDivisor)
	fmt.Println("Press Enter to continue...")
	fmt.Scanln()
}

func clearTerminal() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func menuView() {
	clearTerminal()
	fmt.Println(ASCII)
	fmt.Print(menu)
}

func exitView() {
	fmt.Println("Exiting...")
	clearTerminal()
	os.Exit(0)
}
