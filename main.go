package main

import (
	"bytes"
	"fmt"
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
	"gopkg.in/yaml.v3"
)

const (
	PostsPathMatch    string = "posts/*"
	outputDir         string = "public"
	PostTemplateFile  string = "post.html"
	IndexTemplateFile string = "index.html"
)

type Config struct {
	Author      string `yaml:"author"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Mail        string `yaml:"mail"`
	Github      string `yaml:"github"`
}

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

func parseConfigYaml() {
	var setting Config
	config, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Printf("Error in read config.yaml, %s", err)
	}
	yaml.Unmarshal(config, &setting)
}

func main() {
	http.HandleFunc("/", handlerequest)
	http.ListenAndServe(":8200", nil)
}
