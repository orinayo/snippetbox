package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(resWriter http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(resWriter, req)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	templateSet, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(resWriter, "Internal Server Error", 500)
		return
	}

	err = templateSet.Execute(resWriter, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(resWriter, "Internal Server Error", 500)
	}
}

func showSnippet(resWriter http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(resWriter, req)
		return
	}
	fmt.Fprintf(resWriter, "Display a specific snippet with ID %d", id)
}

func createSnippet(resWriter http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		resWriter.Header().Set("Allow", http.MethodPost)
		http.Error(resWriter, "Method not allowed", 405)
		return
	}
	resWriter.Write([]byte("Create a new snippet"))
}
