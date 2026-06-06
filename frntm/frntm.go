package frntm

import (
	"strings"
	"log"

	"github.com/adrg/frontmatter"
)

func Parse(file []byte) []byte {
	var matter struct {
		Name string   `yaml:"name"`
		Tags []string `yaml:"tags"`
	}

	rest, err := frontmatter.Parse(strings.NewReader(string(file)), &matter)
	if err != nil {
		log.Fatal(err)
	}
	return rest
}
