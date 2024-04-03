package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
)

func main() {
	// create a flag for command
	fileFlag := flag.String("file", "latest-post.txt", "The name of the file")
	flag.Parse()
	//fmt.Println("flag has value", *fileFlag)

	// create new blank HTML file
	file, err := os.Create("first-post.html")
	if err != nil {
		fmt.Println("Error creating HTML file", err)
	}
	// Open the file
	contents, err := os.ReadFile(*fileFlag)
	if err != nil {
		fmt.Print(err)
	}
	//fmt.Println(string(contents))

	// Parse the template
	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// define data structure to pass to the template
	// capitalized to be accessed in tmpl
	type Content struct {
		Content string
	}

	// create an instance of Content
	contentData := Content{
		Content: string(contents),
	}

	// execute the template with the content data and output to the HTML file.
	if err := tmpl.Execute(file, contentData); err != nil {
		fmt.Println("Error executing template:", err)
		return
	}
}
