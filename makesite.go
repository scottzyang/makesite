package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"strings"
)

func main() {
	// declare variables and structs
	var textFile *os.File
	var contents []byte
	var fileContent []byte

	// define data structure to pass to the template
	// capitalized to be accessed in tmpl
	type Content struct {
		Content string
	}

	// create a flag for command
	fileFlag := flag.String("file", "latest-post.txt", "The name of the file")
	flag.Parse()
	//fmt.Println("flag has value", *fileFlag)

	// grab the file name from the command line
	fileName := strings.TrimSuffix(*fileFlag, ".txt")

	// create new blank HTML file with the same name as the text file
	file := createNewHtmlFile(fileName)

	// Double check if the file exists, if not then create it
	fileContent = createOrOpenTextFile(fileFlag, textFile, contents, fileContent)

	// Parse the template
	tmpl := parseTemplate()

	// create an instance of Content
	contentData := Content{
		Content: string(fileContent),
	}

	// execute the template with the content data and output to the HTML file.
	execute(tmpl, file, contentData)
}

func checkFileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false // File does not exist
	}
	return true // File exists
}

func createNewHtmlFile(fileName string) *os.File {
	file, err := os.Create(fileName + ".html")
	if err != nil {
		fmt.Println("Error creating HTML file", err)
	}

	return file
}

func createOrOpenTextFile(fileFlag *string, textFile *os.File, contents []byte, fileContent []byte) []byte {
	// verify that text file exists, if not then create it
	if !checkFileExists(*fileFlag) {
		var err error
		textFile, err = os.Create(*fileFlag)
		if err != nil {
			fmt.Println("Error creating text file", err)
		}
	} else {
		// Open the file
		var err error
		contents, err = os.ReadFile(*fileFlag)
		if err != nil {
			fmt.Print(err)
		}
	}

	// if the text file is nil, then there is no content yet
	if textFile != nil {
		fileContent = []byte("No content yet")
	} else {
		fileContent = contents
	}

	return fileContent
}

func parseTemplate() *template.Template {
	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		fmt.Println("Error parsing template:", err)
	}

	return tmpl
}

func execute(tmpl *template.Template, file *os.File, contentData interface{}) {
	if err := tmpl.Execute(file, contentData); err != nil {
		fmt.Println("Error executing template:", err)
		return
	}
}
