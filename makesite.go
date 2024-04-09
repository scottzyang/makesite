package main

import (
	"flag"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() {
	// initiate measurement of execution time
	start := time.Now()
	// Parse the template
	tmpl := parseTemplate()

	// create new flags
	fileFlag := flag.String("file", "latest-post.txt", "The name of the file.")
	dirFlag := flag.String("dir", "text", "Directory to search for text files.")
	flag.Parse()

	// create single post
	createSinglePost(tmpl, *fileFlag)

	// create multiple posts
	createMultiplePosts(tmpl, *dirFlag)
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}

func createMultiplePosts(tmpl *template.Template, dirFlag string) {
	// total amount of posts
	var totalPosts int
	var postSize float64

	// create multiple posts
	err := filepath.Walk(dirFlag, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error traversing directory", err)
		}
		if path != dirFlag {
			tmpFileFlag := strings.TrimPrefix(path, "text/")
			createSinglePost(tmpl, tmpFileFlag)

			tmpFileFlag = strings.TrimSuffix(tmpFileFlag, ".txt")

			// calculate size
			fileInfo, err := os.Stat("posts/" + tmpFileFlag + ".html")
			if err != nil {
				fmt.Println(err)
			}
			fileSizeKB := float64(fileInfo.Size()) / 1024.0
			fileSizeKB = math.Round(fileSizeKB*10) / 10

			postSize += float64(fileSizeKB)

			totalPosts++
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error initiating directory traversal", err)
	}

	successMessage(totalPosts, postSize)
}

func successMessage(totalPosts int, postSize float64) {
	// style the output
	green := color.New(color.FgGreen)
	boldGreen := color.New(color.FgGreen).Add(color.Bold)
	green.Printf("Success! Generated ")
	boldGreen.Printf("%d", totalPosts)
	green.Printf(" posts. (%.1f kB) \n", postSize)
}

func createSinglePost(tmpl *template.Template, fileFlag string) {
	var fileContent []byte

	// new struct to pass into tmpl
	type Content struct {
		Text string
	}

	// create new blank HTML file with the same name as the text file
	htmlFile := createNewHtmlFile(fileFlag)

	// Double check if the file exists, if not then create it
	fileContent = createOrReadTextFile(fileFlag)

	// add content to instance of Content
	contentData := Content{
		Text: string(fileContent),
	}

	// take parsed template, with the contentData and put it in HTML file
	execute(tmpl, htmlFile, contentData)
}

func fileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false // File does not exist
	}
	return true // File exists
}

func createNewHtmlFile(fileFlag string) *os.File {
	fileName := strings.TrimSuffix(fileFlag, ".txt")

	file, err := os.Create("./posts/" + fileName + ".html")
	if err != nil {
		fmt.Println("Error creating HTML file", err)
	}

	return file
}

func createOrReadTextFile(fileFlag string) []byte {
	var contents []byte
	var fileContent []byte
	var textFile *os.File
	var newTextContent []byte

	// Verify that text file exists, if not then create it
	if fileExists("./text/" + fileFlag) {

		var err error
		contents, err = os.ReadFile("./text/" + fileFlag)
		if err != nil {
			fmt.Print(err)
		}
	} else {
		var err error
		textFile, err = os.Create("./text/" + fileFlag)
		fmt.Println("Creating textfile", textFile)
		if err != nil {
			fmt.Println("Error creating text file", err)
		}

		emptyContent := []byte("No content yet!")
		_, err = textFile.Write(emptyContent)
		fmt.Println("Writing textfile")
		if err != nil {
			fmt.Println("Error writing to new text file", err)
		}

		newTextContent, err = os.ReadFile("./text/" + fileFlag)
		fmt.Println("Reading newly created text file", fileFlag)
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
