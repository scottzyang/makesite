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
	var contents []byte
	var fileContent []byte

	// define data structure to pass to the template
	// capitalized to be accessed in tmpl
	type Content struct {
		Content string
	}

	// create new flags
	fileFlag := flag.String("file", "latest-post.txt", "The name of the file")
	flag.Parse()

	// grab the file name from the command line
	fileName := strings.TrimSuffix(*fileFlag, ".txt")

	// create new blank HTML file with the same name as the text file
	file := createNewHtmlFile(fileName)

	// Double check if the file exists, if not then create it
	fileContent = createOrReadTextFile(fileFlag, contents)

	// Parse the template
	tmpl := parseTemplate()

	// create an instance of Content
	contentData := Content{
		Content: string(fileContent),
	}

	// execute the template with the content data and output to the HTML file.
	execute(tmpl, file, contentData)
}

func fileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false // File does not exist
	}
	return true // File exists
}

func createNewHtmlFile(fileName string) *os.File {
	file, err := os.Create("./posts/" + fileName + ".html")
	if err != nil {
		fmt.Println("Error creating HTML file", err)
	}

	return file
}

func createOrReadTextFile(fileFlag *string, contents []byte) []byte {
	var fileContent []byte
	var textFile *os.File
	var newTextContent []byte

	// Verify that text file exists, if not then create it
	if fileExists("./text/" + *fileFlag) {

		var err error
		contents, err = os.ReadFile("./text/" + *fileFlag)
		fmt.Println("Reading textfile:")
		if err != nil {
			fmt.Print(err)
		}
	} else {
		var err error
		textFile, err = os.Create("./text/" + *fileFlag)
		fmt.Println("Creating textfile:")
		if err != nil {
			fmt.Println("Error creating text file", err)
		}

		emptyContent := []byte("No content yet!")
		_, err = textFile.Write(emptyContent)
		fmt.Println("Writing textfile:")
		if err != nil {
			fmt.Println("Error writing to new text file", err)
		}

		newTextContent, err = os.ReadFile("./text/" + *fileFlag)
		fmt.Println("Reading newly created text file:")
		if err != nil {
			fmt.Println("Error reading newly created text file", err)
		}
	}

	// If the text file is nil, then there is no content yet
	if textFile != nil {
		fileContent = newTextContent
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
