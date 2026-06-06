package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"github.com/nielmin/go-ssg/frntm"
)

func mdToHTML(md []byte) []byte {
	// parse frontmatter before rendering
	fmless := frntm.Parse(md)

	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(fmless)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// sanitize policy
	s := bluemonday.UGCPolicy()

	// sanitize html
	return s.SanitizeBytes(markdown.Render(doc, renderer))
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
	mux := http.NewServeMux()

	file, err := os.ReadFile("test.md")
	if err != nil {
		log.Fatal(err)
	}

	html := mdToHTML(file)

	newFile := os.WriteFile("index.html", html, 0o666)
	if newFile != nil {
		log.Fatal(newFile)
	}

	mux.HandleFunc("/", htmlRead)

	log.Print("Listening on port " + port + "...")
	http.ListenAndServe(":"+port+"", mux)
}
