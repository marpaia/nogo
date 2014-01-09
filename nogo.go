package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
  nogo new
  nogo new [topic]
  nogo new [topic] [event]
  nogo ls
  nogo ls [topic]
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

func generateTimestamp() string {
	layout := "2006-1-2_"
	timeNow := time.Now()
	return timeNow.Format(layout)
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
	fileName := fmt.Sprintf("%s%s.md", generateTimestamp(), event)
	fullPath := filepath.Join(notesDirectory, fileName)

	// Create the new notes file
	createFile(fullPath)

	// Open the new file in your editor of choice
	openFile(fullPath)
}

func parseFilename(filename string) string {
	dateName := strings.Split(filename, "_")
	if len(dateName) == 2 {
		nameExt := strings.SplitAfter(dateName[1], ".")
		if len(nameExt) >= 2 {
			filenameDashes := nameExt[0 : len(nameExt)-1]
			parsedFilename := strings.Join(filenameDashes, "")
			parsedFilename = strings.Replace(parsedFilename, "-", " ", -1)
			parsedFilename = strings.TrimRight(parsedFilename, ".")
			return fmt.Sprintf("%s (%s)", parsedFilename, dateName[0])
		}
	}

	return "unparsable filename"
}

func listTopics() {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println("Oops, couldn't read that directory:", directory)
		os.Exit(1)
	}

	fmt.Println()
	if len(files) == 0 {
		fmt.Println("looks like there aren't any topcis to list!")
	} else {
		fmt.Println("all topics:")
	}

	for _, file := range files {
		fmt.Println("  ", file.Name())
	}
	fmt.Println()
}

func listNotes(topic string) {
	topicDir := fmt.Sprintf("%s/%s", directory, topic)
	files, err := ioutil.ReadDir(topicDir)
	if err != nil {
		fmt.Println("Oops, couldn't read that directory:", topicDir)
		os.Exit(1)
	}

	fmt.Println()
	if len(files) == 0 {
		fmt.Println("looks like there aren't any notes to list!")
	} else {
		fmt.Printf("notes in %s:\n", topic)
	}

	for _, file := range files {
		fmt.Println("  ", parseFilename(file.Name()))
	}

	fmt.Println()

}

func editFile(topic, target string) {
	topicDir := fmt.Sprintf("%s/%s", directory, topic)
	files, err := ioutil.ReadDir(topicDir)
	if err != nil {
		fmt.Println("Oops, couldn't read that directory:", topicDir)
	}

	target = strings.Replace(target, " ", "-", -1)
	var found bool

	for _, file := range files {
		if strings.Contains(file.Name(), target) {
			found = true
			fullPath := fmt.Sprintf("%s/%s/%s", directory, topic, file.Name())
			openFile(fullPath)
		}
	}

	if !found {
		fmt.Println("Oops, couldn't find that file!")
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) == 1 {
		help(0)
	}
	action := os.Args[1]

	switch action {
	case "new":
		if len(os.Args) == 2 {
			handleNew()
		} else if len(os.Args) == 3 {
			handleNewTopic(os.Args[2])
		} else if len(os.Args) == 4 {
			handleNewEvent(os.Args[2], os.Args[3])
		} else {
			help(1)
		}
	case "ls":
		if len(os.Args) == 2 {
			listTopics()
		} else if len(os.Args) == 3 {
			listNotes(os.Args[2])
		} else {
			help(1)
		}
	case "edit":
		if len(os.Args) < 4 {
			fmt.Println("I need a topic and a note to edit!")
			os.Exit(1)
		}
		editFile(os.Args[2], strings.Join(os.Args[3:], " "))
	default:
		help(0)
	}
}
