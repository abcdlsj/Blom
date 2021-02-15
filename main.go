package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	PostsPathMatch    string = "posts/*"
	outputDir         string = "public"
	PostTemplateFile  string = "post.html"
	IndexTemplateFile string = "index.html"
)

type Post struct {
	Title   string
	Date    string
	Summary string
	Body    string
	File    string
}

func handlerequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "" {
		posts := getPosts()
		t := template.New(IndexTemplateFile)
		t, _ = t.ParseFiles(IndexTemplateFile)
		t.Execute(w, posts)
	} else {
		f := "posts/" + r.URL.Path[1:] + ".md"
		fileread, _ := ioutil.ReadFile(f)
		lines := strings.Split(string(fileread), "\n")
		title := string(lines[1])
		date := string(lines[2])
		summary := string(lines[3])
		body := strings.Join(lines[5:], "\n")
		htmlByte, err := markdown2HTML([]byte(body))
		if err != nil {
			log.Fatal("markdown2HTML error!")
		}
		post := Post{title, date, summary, string(htmlByte), r.URL.Path[1:]}
		t := template.New(PostTemplateFile)
		t, _ = t.ParseFiles(PostTemplateFile)
		t.Execute(w, post)
	}
}
func markdown2HTML(src []byte) ([]byte, error) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithHardWraps()),
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"))))

	var buf bytes.Buffer
	if err := markdown.Convert(src, &buf); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func getPosts() []Post {
	a := []Post{}
	files, _ := filepath.Glob("posts/*")
	for _, f := range files {
		file := strings.Replace(f, "posts/", "", -1)
		file = strings.Replace(file, ".md", "", -1)
		fileread, _ := ioutil.ReadFile(f)
		lines := strings.Split(string(fileread), "\n")
		title := string(lines[1])
		date := string(lines[2])
		summary := string(lines[3])
		body := strings.Join(lines[5:], "\n")
		htmlByte, err := markdown2HTML([]byte(body))
		if err != nil {
			log.Fatal("markdown2HTML error!")
		}
		a = append(a, Post{title, date, summary, string(htmlByte), file})
	}
	return a
}

func main() {
	http.HandleFunc("/", handlerequest)
	http.ListenAndServe(":8100", nil)
}
