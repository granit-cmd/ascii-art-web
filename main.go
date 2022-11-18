package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	ascii_art "ascii-art-web/ascii-art"
)

var output string

func DrawWeb(writer http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		http.Error(writer, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	content, err := os.ReadFile("index.html")
	if err != nil {
		return
	}

	strContent := string(content)

	t, err := template.New("ascii-art").Parse(strContent)
	if err != nil {
		http.Error(writer, "File not found: index.html", http.StatusInternalServerError)
		return
	}

	text, textErr := req.PostForm["text"]
	font, fontErr := req.PostForm["font"]

	if !fontErr || !textErr {
		http.Error(writer, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	fontStr := font[0]
	textStr := text[0]
	output, err = ascii_art.Draw(textStr, fontStr)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	err = t.Execute(writer, output)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func StartPage(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if req.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("./index.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "File not found: index.html", http.StatusInternalServerError)
	}
	err = t.Execute(w, output)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func main() {
	// http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/ascii-art", DrawWeb)
	http.HandleFunc("/", StartPage)
	fmt.Println("Started server at http://localhost:4000/ ")
	if err := http.ListenAndServe(":4000", nil); err != nil {
		fmt.Println(err)
		return
	}
}
