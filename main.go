package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gernest/front"
	"github.com/russross/blackfriday"
)

type Post struct {
	Title   string
	Date    string
	Summary string
	Body    string
	File    string
}

func handlerequest(w http.ResponseWriter, r *http.Request) {
	//var err error
	var post Post
	if r.URL.Path[1:] == "" {
		posts := getPosts()
		t := template.New("index.html")
		t, _ = t.ParseFiles("templates/index.html")
		t.Execute(w, posts)
	} else {
		fpath := "posts/" + r.URL.Path[1:] + ".md"
		fileread, _ := ioutil.ReadFile(fpath)

		m := front.NewMatter()
		m.Handle("---", front.YAMLHandler)
		f, body, err := m.Parse(strings.NewReader(string(fileread)))
		title := fmt.Sprintf("%v", f["title"])
		date := fmt.Sprintf("%v", f["date"])
		summary := fmt.Sprintf("%v", f["summary"])
		body = string(blackfriday.MarkdownCommon([]byte(body)))
		post = Post{title, date, summary, body, r.URL.Path[1:]}
		t := template.New("post.html")
		t, _ = t.ParseFiles("templates/post.html")
		err = t.Execute(w, post)
		if err != nil {
			fmt.Println("error", err)
		}
	}
}

func getPosts() []Post {
	a := []Post{}
	files, _ := filepath.Glob("posts/*")
	for _, f := range files {
		file := strings.Replace(f, "posts/", "", -1)
		file = strings.Replace(file, ".md", "", -1)
		fileread, _ := ioutil.ReadFile(f)
		m := front.NewMatter()
		m.Handle("---", front.YAMLHandler)
		f, body, err := m.Parse(strings.NewReader(string(fileread)))
		if err != nil {
			fmt.Println("error", err)
		}
		title := fmt.Sprintf("%v", f["title"])
		date := fmt.Sprintf("%v", f["date"])
		summary := fmt.Sprintf("%v", f["summary"])
		body = string(blackfriday.MarkdownCommon([]byte(body)))
		a = append(a, Post{title, date, summary, body, file})
	}
	return a
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", handlerequest)
	http.ListenAndServe(":8080", nil)
}
