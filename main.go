package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

func ensureTodoDir() string {
	dir := "Note"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal("❌ Failed to create Note directory:", err)
		}
	}
	return dir
}
func directory(title string) {
	dir := ensureTodoDir()
	filename := dir + "/" + title + ".txt"
	ReadWriteFile(filename)
}

func ReadWriteFile(filename string) string {
	fmt.Println("write down here....")
	fmt.Print(" › ")
	read := bufio.NewScanner(os.Stdin)

	for {
		if !read.Scan() {
			break
		}
		exits := strings.TrimSpace(read.Text())

		if exits == "" {
			// instead of breaking, ask again
			fmt.Println("❌ Input cannot be empty. Please try again.")
			fmt.Print(" › ")
			continue
		}

		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if err := read.Err(); err != nil {
			log.Fatal(err)
		}
		if _, err := file.WriteString(exits + "\n"); err != nil {
			log.Fatal(err)
		}
		file.Close()
		break
	}

	fmt.Println("✅ Saved to file:", filename)
	Menu()
	MenuList()
	return filename
}

func EditWriteFile(filename string) {
	dir := ensureTodoDir()
	filePath := dir + "/" + filename + ".txt"

	_, data := os.ReadFile(filePath)
	if data != nil {
		fmt.Println("❌ Cannot find the file:", filePath)
		return
	}
	cmd := exec.Command("nano", filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ File updated:", filename)
	Menu()
	MenuList()

}

// ANSI color codes
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Cyan    = "\033[36m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	White   = "\033[37m"
)

func Welcome() {
	fmt.Println(Green, "     ◆  Welcome to WriteDown  ◆      ", Green)
	fmt.Println(Reset, "", Reset)
	fmt.Println("This is a simple Cli Note Taking App ")
	fmt.Println("")
	selection()
	fmt.Println("|     " + Yellow + "/help   Show help" + Reset + "                  |")
	fmt.Println("|     " + Yellow + "/about  About this app" + Reset + "             |")
	fmt.Println("")
}

func Menu() {

	fmt.Println("————————————— WriteDown ———————————————")
	fmt.Println("|                                        |")
	fmt.Println("|              1. Add                    |")
	fmt.Println("|              2. Edit                   |")
	fmt.Println("|              3. Delete                 |")
	fmt.Println("|              4. list                   |")
	fmt.Println("|              5. Exit                   |")
	fmt.Println("|                                        |")
	fmt.Println("———————————————————————————————————————")
	fmt.Println(" ")
}
func MenuList() {
	var ch int
	fmt.Print(" › ")
	fmt.Scan(&ch)
	bufio.NewReader(os.Stdin).ReadString('\n')
	switch ch {
	case 1:
		fmt.Println(" ")
		fmt.Println("———————————————————————————————————————")
		fmt.Println("|                                     |")
		fmt.Println("|      Give a Name to the topic       |")
		fmt.Println("|                                     |")
		fmt.Println("———————————————————————————————————————")
		fmt.Print(" › ")
		for {
			read := bufio.NewScanner(os.Stdin)
			read.Scan()
			fmt.Println(" ")
			input := strings.TrimSpace(read.Text())
			if input == "" {
				fmt.Println("❌ Input cannot be empty. Please try again.")
				fmt.Print(" › ")
				continue
			}
			directory(input)
			break
		}

	case 2:
		fmt.Println("———————————————————————————————————————")
		fmt.Println("|                                     |")
		fmt.Println("|                Edit                 |")
		fmt.Println("|                                     |")
		fmt.Println("———————————————————————————————————————")
		fmt.Println("")
		dir := ensureTodoDir()
		cmd, err := exec.Command("ls", dir).Output()
		if err != nil {
			log.Fatal(err)
		}
		file := strings.Split(strings.TrimSpace(string(cmd)), "\n")
		for i, files := range file {
			fmt.Println(i+1, "→", files)
		}
		fmt.Println("🔍 Enter the name of task ")
		fmt.Print(" › ")
		for {
			edit := bufio.NewScanner(os.Stdin)
			edit.Scan()
			input := strings.TrimSpace(edit.Text())
			if input == "" {
				fmt.Println("❌ Input cannot be empty. Please try again.")
				fmt.Print(" › ")
				continue
			}
			if _, err := os.Stat("Todo/" + input + ".txt"); os.IsNotExist(err) {
				fmt.Printf("❌ File '%s' does not exist. Do you want to create it? (y/n): ", input)
				choice := bufio.NewScanner(os.Stdin)
				choice.Scan()
				answer := strings.ToLower(strings.TrimSpace(choice.Text()))
				if answer == "y" || answer == "yes" {
					directory(input) // create the file and start writing
				} else if answer == "n" || answer == "N" {
					Menu()
					MenuList()
				}
				break
			}
			EditWriteFile(input)
			break
		}
	case 3:
		dir := ensureTodoDir()
		cmd, err := exec.Command("ls", dir).Output()
		if err != nil {
			log.Fatal(err)
		}
		files := strings.Split(strings.TrimSpace(string(cmd)), "\n")
		if len(files) == 0 || (len(files) == 1 && files[0] == "") {
			fmt.Print("Do you want to create a new note? (y/n): ")
			choice := bufio.NewScanner(os.Stdin)
			choice.Scan()
			answer := strings.ToLower(strings.TrimSpace(choice.Text()))
			if answer == "y" || answer == "yes" {
				fmt.Println("———————————————————————————————————————")
				fmt.Println("|                                     |")
				fmt.Println("|      Give a Name to the topic       |")
				fmt.Println("|                                     |")
				fmt.Println("———————————————————————————————————————")
				fmt.Print(" › ")
				for {
					read := bufio.NewScanner(os.Stdin)
					read.Scan()
					input := strings.TrimSpace(read.Text())
					if input == "" {
						fmt.Println("❌ Input cannot be empty. Please try again.")
						fmt.Print(" › ")
						continue
					}
					directory(input)
					break
				}
			} else {
				fmt.Println("Returning to menu...")
				Menu()
				MenuList()
			}
			return
		}
		// Display the existing files
		for i, file := range files {
			fmt.Println(i+1, "→", file)
		}

		for {
			fmt.Println("🔍 Enter the name of task to delete")
			fmt.Print(" › ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())

			if input == "" {
				fmt.Println("❌ Input cannot be empty. Please try again.")
				Menu()
				MenuList()
			}

			filePath := dir + "/" + input + ".txt"
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				fmt.Println("❌ File does not exist. Please try again.")
				continue
			}

			// Confirm deletion
			fmt.Printf("⚠ Are you sure you want to delete '%s'? (y/n): ", input)
			confirm := bufio.NewScanner(os.Stdin)
			confirm.Scan()
			answer := strings.ToLower(strings.TrimSpace(confirm.Text()))
			if answer == "y" || answer == "yes" {
				deleteFile(input)
			} else {
				fmt.Println("❌ Deletion cancelled.")
				Menu()
				MenuList()
			}
			break
		}

	case 4:
		List()
	case 5:
		fmt.Println("Have a good day!")
	default:
		fmt.Println("something wrong ! Try again")
		Menu()
		MenuList()
	}
}
func deleteFile(filename string) {
	dir := ensureTodoDir()
	filePath := dir + "/" + filename + ".txt"

	_, data := os.ReadFile(filePath)
	if data != nil {
		fmt.Println("❌ Cannot find the file:", filePath)
		return
	}
	del := os.Remove(filePath)
	if del != nil {
		log.Fatal("❌ Error deleting file:", del)
	}

	fmt.Println("✅ File deleted:", filename+".txt")
	Menu()
	MenuList()
}
func List() {
	dir := ensureTodoDir()
	cmd, err := exec.Command("ls", dir).Output()
	if err != nil {
		log.Fatal(err)
	}

	bar := "▮ "
	for i := 1; i <= len(bar); i++ {
		fmt.Printf("\r%s", bar[:i])
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("")
	file := strings.Split(strings.TrimSpace(string(cmd)), "\n")
	for i, files := range file {
		fmt.Println(i+1, "→", files)
	}

	fmt.Println("Done!")
	time.Sleep(1000 * time.Millisecond)

	Menu()
	MenuList()
}
func clearTerminal() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default: // Linux or macOS
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
func loadingBar(steps int) {
	fmt.Print(Cyan + "Loading: " + Reset)
	for i := 0; i < steps; i++ {
		fmt.Print(Green + "▮" + Reset)
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println(Green + " ✅" + Reset)
	fmt.Println("")
}
func About() {
	fmt.Println(Cyan + "================ Writedown v1.0 ================" + Reset)
	fmt.Println(Green + "A simple, terminal-based note-taking application written in Go." + Reset)
	fmt.Println("")
	fmt.Println("Features:")
	fmt.Println(" ⇨ Add notes with a title and description")
	fmt.Println(" ⇨ Edit existing notes using the nano editor")
	fmt.Println(" ⇨ List all saved notes with a visual loading indicator")
	fmt.Println(" ⇨ Delete notes safely with confirmation")
	fmt.Println(" ⇨ Easy navigation with interactive menus or number-based selection")
	fmt.Println(" ⇨ Cross-platform support for Windows, Linux, and macOS terminals")
	fmt.Println("")
	fmt.Println("Created by: Aswanidev VS")
	fmt.Println(Cyan + "===============================================" + Reset)
	fmt.Println("\nPress Enter to return to the main menu...")
	// fmt.Scanln() // waits for Enter
}

func selection() string {
	items := []string{"Start", "About", "Exit"}

	prompt := promptui.Select{
		Label: "\nselect\n",
		Items: items,
		Size:  len(items),
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	// fmt.Printf("You selected %q\n", result)
	return result
}

func startup() {
	clearTerminal()
	fmt.Println(Green, "      ◆           Welcome to WriteDown           ◆", Reset)
	fmt.Println("")
	fmt.Println("         WriteDown is a simple, terminal-based")
	fmt.Println("         note-taking application written in Go.")
	fmt.Println("         It allows users to quickly create, edit,")
	fmt.Println("         list, and delete notes directly from the")
	fmt.Println("         command line, making it lightweight and")
	fmt.Println("         fast without requiring any heavy GUI.")
	fmt.Println("")
	fmt.Println(Green, "      ◆——————————————————————————————————————————◆", Reset)

	choice := selection()

	switch choice {
	case "Start":
		fmt.Println("🚀 Starting the app...")
		loadingBar(5)
		clearTerminal()
		Menu()
		MenuList()
	case "About":
		About()
		fmt.Scanln()
		startup()
	case "Exit":
		clearTerminal()
	default:
		fmt.Println("❌ Unknown selection")
	}
}

func main() {
	startup()
}
