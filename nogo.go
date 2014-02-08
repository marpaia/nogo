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

// the help text that gets displayed when something goes wrong or when you run
// help
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

// directory holds the string that represents the notes that your notes are
// stored in
var directory string

// editor holds the string that represents the editor that you use to edit your
// notes. it defaults to "vim" but it is overriden by the $EDITOR environment
// variable
var editor string

// init gathers the environment variables on the system and interpolates
// relevant strings
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

// help prints the help text to stdout and then exits with a given exit code
func help(exit int) {
	fmt.Println(helpText)
	os.Exit(exit)
}

// acceptInput asks the command-line user a given question and returns the text
// that they submitted
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

// normalizeString accepts a pointer to a string and modifies it to replace all
// spaces with dashes
func normalizeString(s *string) {
	*s = strings.Replace(*s, " ", "-", -1)
}

// openFile accepts a file name as a parameter and opens it with the editor
// that is stored in the "editor" variable
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

// createFile creates a given filename if it doesn't exist
func createFile(fileName string) {
	_, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Oops, I couldn't create the file:", fileName)
	}
}

// createDir creates a directory if the directory path doesn't already exist
func createDir(dirName string) {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		fmt.Println("Oops, I couldn't make the directory:", dirName)
		os.Exit(1)
	}
}

// handleNew is the function that gets executed when you type "nogo new"
func handleNew() {
	// Ask for the topic of notes
	topic := acceptInput("Enter the notes topic: ")
	normalizeString(&topic)

	handleNewTopic(topic)
}

// handleNewTopic is the function that gets executed when you type
// "nogo new [topic]"
func handleNewTopic(topic string) {
	// Ask for the event name of the notes
	event := acceptInput("Enter the event name: ")
	normalizeString(&event)

	handleNewEvent(topic, event)
}

// handleNewEvent is the function that gets executed when you type
// "nogo new [topic] [event]"
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

// parseFilename accepts a string as input and returns a new string, in which,
// all dashes are replaced with spaces and the file extension is stripped
func parseFilename(filename string) string {
	return strings.Replace(strings.TrimSuffix(filename, ".md"), "-", " ", -1)
}

// isDir returns true if the given path is a directory
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

// listTopics lists the existing topics, while avoiding .git artifacts
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
		fmt.Println("looks like there aren't any topics to list!")
	} else {
		fmt.Println("all topics:")
	}

	for _, file := range toPrint {
		fmt.Println("  ", file)
	}
	fmt.Println()
}

// findTopicBySubstring accepts a substring of an existing topic and returns
// the complete topic as well as an error indicating that either something went
// wrong or the topic substring wasn't found
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

// listNotes accepts a topic substring and lists all of the notes in that topic
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

// findFileInList takes a list of files, a target filename substring and the
// relevant topic. it then iterates through the files, identifies the proper
// file and opens it in your editor of choice
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

// editFileWithoutTopic is the function that gets called when you type
// "nogo edit"
func editFileWithoutTopic() {
	topic := acceptInput("\nWhat topic would you like to edit? ")
	editFileInTopic(topic)
}

func findFilesInTopic(topic string) ([]os.FileInfo, string) {
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
	return files, completeTopic
}

// editFileInTopic is the function that gets called when you type
// "nogo edit [topic]". it also gets called as the second part of
// "nogo edit" once the topic is given by the user.
func editFileInTopic(topic string) {
	files, completeTopic := findFilesInTopic(topic)
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

// editFile accepts a topic substring and a target file name substring and
// opens it in your editor of choice
func editFile(topic, target string) {
	files, completeTopic := findFilesInTopic(topic)
	findFileInList(files, target, completeTopic)
}

// main is what gets executed when you run nogo from the command line
func main() {
	if len(os.Args) == 1 {
		help(0)
	}
	action := os.Args[1]
	commandArgs := len(os.Args) - 2

	switch action {
	case "help":
		help(0)
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
		help(1)
	}
}
