package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// this will store notes in ~/notes if the NOGODIR environment variable isn't
// set
var NotesSubDir = "notes"

// this will be used as the default editor if the EDITOR environment variable
// isn't set
var DefaultEditor = "vim"

var helpText = `
nogo - the notes helper

actions:

  nogo help

  nogo new
  nogo new [topic]
  nogo new [topic] [event]

  nogo ls
  nogo ls [topic]

  nogo edit
  nogo edit [topic]
  nogo edit [topic] [note name substring]
`

var directory string
var editor string

func init() {
	directory = os.Getenv("NOGODIR")
	if directory == "" {
		directory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), NotesSubDir)
	}

	editor = os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}
}

func help(exit int) {
	fmt.Println(helpText)
	os.Exit(exit)
}

func acceptInput(question string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(question)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Oops, what was that?")
		os.Exit(1)
	}

	return strings.TrimSpace(response)
}

func normalizeString(s *string) {
	*s = strings.Replace(*s, " ", "-", -1)
}

func openFile(fileName string) {
	cmd := exec.Command(editor, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Couldn't open the file:", err)
		os.Exit(1)
	}
}

func createFile(fileName string) {
	_, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Oops, I couldn't create the file:", fileName)
	}
}

func createDir(dirName string) {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		fmt.Println("Oops, I couldn't make the directory:", dirName)
		os.Exit(1)
	}
}

func gatherFiles() {
	files := []string{}
	visit := func(path string, f os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	}
	filepath.Walk(directory, visit)
	fmt.Println(files)

}

func handleNew() {
	// Ask for the topic of notes
	topic := acceptInput("Enter the notes topic: ")
	normalizeString(&topic)

	handleNewTopic(topic)
}

func handleNewTopic(topic string) {
	// Ask for the event name of the notes
	event := acceptInput("Enter the event name: ")
	normalizeString(&event)

	handleNewEvent(topic, event)
}

func handleNewEvent(topic, event string) {
	// Create the topic folder if it doesn't exist
	notesDirectory := filepath.Join(directory, topic)
	createDir(notesDirectory)

	// Make the filename and filepath for the new notes file
	fileName := fmt.Sprintf("%s.md", event)
	fullPath := filepath.Join(notesDirectory, fileName)

	// Create the new notes file
	createFile(fullPath)

	// Open the new file in your editor of choice
	openFile(fullPath)
}

func parseFilename(filename string) string {
	return strings.Replace(strings.TrimSuffix(filename, ".md"), "-", " ", -1)
}

func isDir(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	default:
		return false
	}

}

func listTopics() {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println("Oops, couldn't read that directory:", directory)
		os.Exit(1)
	}

	fmt.Println()

	toPrint := []string{}
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", directory, file.Name())
		if isDir(path) && !strings.HasPrefix(file.Name(), ".git") {
			toPrint = append(toPrint, file.Name())
		}
	}

	if len(toPrint) == 0 {
		fmt.Println("looks like there aren't any topcis to list!")
	} else {
		fmt.Println("all topics:")
	}

	for _, file := range toPrint {
		fmt.Println("  ", file)
	}
	fmt.Println()
}

func findTopicBySubstring(substring string) (string, error) {
	topics, err := ioutil.ReadDir(directory)
	if err != nil {
		return "", err
	}
	for _, topic := range topics {
		if strings.Contains(topic.Name(), substring) {
			return topic.Name(), nil
		}
	}

	return "", errors.New("Couldn't find that filename")
}

func listNotes(topic string) {
	completeTopic, err := findTopicBySubstring(topic)
	if err != nil {
		fmt.Println("Oops, couldn't find that topic:", topic)
		os.Exit(1)
	}
	topicDir := fmt.Sprintf("%s/%s", directory, completeTopic)
	files, err := ioutil.ReadDir(topicDir)
	if err != nil {
		fmt.Println("Oops, couldn't read that directory:", topicDir)
		os.Exit(1)
	}

	fmt.Println()
	if len(files) == 0 {
		fmt.Println("looks like there aren't any notes to list!")
	} else {
		fmt.Printf("notes in %s:\n", completeTopic)
	}

	for _, file := range files {
		fmt.Println("  ", parseFilename(file.Name()))
	}

	fmt.Println()

}

func findFileInList(files []os.FileInfo, target, completeTopic string) {
	target = strings.Replace(target, " ", "-", -1)
	var found bool

	for _, file := range files {
		if strings.Contains(file.Name(), target) {
			found = true
			fullPath := fmt.Sprintf("%s/%s/%s", directory, completeTopic, file.Name())
			openFile(fullPath)
			os.Exit(0)
		}
	}

	if !found {
		fmt.Println("Oops, couldn't find that file!")
		os.Exit(1)
	}
}

func editFileWithoutTopic() {
	topic := acceptInput("\nWhat topic would you like to edit? ")
	editFileInTopic(topic)
}

func editFileInTopic(topic string) {
	completeTopic, err := findTopicBySubstring(topic)
	if err != nil {
		fmt.Println("Oops, couldn't find that topic:", topic)
		os.Exit(1)
	}
	topicDir := fmt.Sprintf("%s/%s", directory, completeTopic)
	files, err := ioutil.ReadDir(topicDir)
	if err != nil {
		fmt.Println("Oops, couldn't read that directory:", topicDir)
		os.Exit(1)
	}
	if len(files) > 0 {
		fmt.Printf("\nFiles in %s:\n", completeTopic)
	}
	for _, file := range files {
		fmt.Printf("   %s\n", parseFilename(file.Name()))
	}
	fmt.Println()
	file := acceptInput(fmt.Sprintf("What file would you like to edit in %s? ", completeTopic))
	findFileInList(files, file, completeTopic)
}

func editFile(topic, target string) {
	completeTopic, err := findTopicBySubstring(topic)
	if err != nil {
		fmt.Println("Oops, couldn't find that topic:", topic)
		os.Exit(1)
	}
	topicDir := fmt.Sprintf("%s/%s", directory, completeTopic)
	files, err := ioutil.ReadDir(topicDir)
	if err != nil {
		fmt.Println("Oops, couldn't read that directory:", topicDir)
		os.Exit(1)
	}
	findFileInList(files, target, completeTopic)
}

func main() {
	if len(os.Args) == 1 {
		help(0)
	}
	action := os.Args[1]
	commandArgs := len(os.Args) - 2

	switch action {
	case "new":
		switch commandArgs {
		case 0:
			handleNew()
		case 1:
			handleNewTopic(os.Args[2])
		case 2:
			handleNewEvent(os.Args[2], os.Args[3])
		default:
			help(1)
		}
	case "ls":
		switch commandArgs {
		case 0:
			listTopics()
		case 1:
			listNotes(os.Args[2])
		default:
			help(1)
		}
	case "edit":
		switch commandArgs {
		case 0:
			editFileWithoutTopic()
		case 1:
			editFileInTopic(os.Args[2])
		default:
			editFile(os.Args[2], strings.Join(os.Args[3:], "-"))
		}
	default:
		help(0)
	}
}
