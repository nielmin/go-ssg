package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func htmlRead(w http.ResponseWriter, r *http.Request) {
	index, err := os.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(index)
}

func main() {
	port := "8080"

	file, err := os.ReadFile("test.md")
	if err != nil {
		log.Fatal(err)
	}

	var matter struct {
		Name string   `yaml:"name"`
		Tags []string `yaml:"tags"`
	}

	rest, err := frontmatter.Parse(strings.NewReader(string(file)), &matter)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	p := bluemonday.UGCPolicy()
	html := p.SanitizeBytes(mdToHTML(rest))

	newFile := os.WriteFile("index.html", html, 0o666)
	if newFile != nil {
		log.Fatal(newFile)
	}

	mux.HandleFunc("/", htmlRead)

	log.Print("Listening on port " + port + "...")
	http.ListenAndServe(":"+port+"", mux)
}
