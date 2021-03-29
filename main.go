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
)

// no comment
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

type Meta struct {
	Title   string
	Date    string
	Tags    []string
	Summary string
}

type Post struct {
	MetaData Meta
	Body     string
	FileName string
}

type Tag struct {
	Title string
	Link  string
	Count int
}

func createTagPostsMap(posts []*Post) map[string][]*Post {
	result := make(map[string][]*Post)
	for _, post := range posts {
		for _, tag := range post.MetaData.Tags {
			key := strings.ToLower(tag)
			if result[key] == nil {
				result[key] = []*Post{post}
			} else {
				result[key] = append(result[key], post)
			}
		}
	}
	return result
}

func getTagLink(tag string) {
}

func getPostInfo(f string, name string) Post {
	fileread, _ := ioutil.ReadFile(f)
	lines := strings.Split(string(fileread), "\n")
	title := lines[1]
	date := lines[2]
	tags := strings.Split(lines[3], " ")
	summary := lines[4]
	body := strings.Join(lines[6:], "\n")
	htmlByte, err := markdown2HTML([]byte(body))
	if err != nil {
		log.Fatal("markdown2HTML error!")
	}

	return Post{Meta{title, date, tags, summary}, string(htmlByte), name}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "" {
		posts := getPosts()
		t := template.New(IndexTemplateFile)
		t, _ = t.ParseFiles(IndexTemplateFile)
		t.Execute(w, posts)
	} else {
		f := "posts/" + r.URL.Path[1:] + ".md"
		post := getPostInfo(f, r.URL.Path[1:])
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
	var a []Post
	files, _ := filepath.Glob("posts/*")
	for _, f := range files {
		fname := strings.Replace(f, "posts/", "", -1)
		fname = strings.Replace(fname, ".md", "", -1)
		post := getPostInfo(f, fname)
		a = append(a, post)
	}
	return a
}

func main() {
	fmt.Println("now server in localhost:8200")
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8200", nil)
}
