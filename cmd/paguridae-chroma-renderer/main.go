package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"xuejie.space/c/go-quill-chroma-renderer"
)

var filename = flag.String("filename", "", "Filename used to deduct source type")
var outputFile = flag.String("outputFile", "-", "Output file to generate, use '-' to print to stdout")
var readFromFile = flag.Bool("readFromFile", false, "True to read from filename, false to read from stdin")
var styleName = flag.String("style", "algol", "Chroma style to use")

func main() {
	flag.Parse()
	if len(*filename) == 0 {
		log.Fatalf("Filename is required!")
	}

	var err error
	reader := os.Stdin
	if *readFromFile {
		reader, err = os.Open(*filename)
		if err != nil {
			log.Fatal(err)
		}
	}

	contentBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	content := string(contentBytes)

	lexer := lexers.Match(*filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)
	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		log.Fatal(err)
	}

	writer := os.Stdout
	if *outputFile != "-" {
		writer, err = os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	style := styles.Get(*styleName)
	if style == nil {
		style = styles.Fallback
	}

	err = renderer.FormatToChroma(writer, style, iterator)
	if err != nil {
		log.Fatal(err)
	}
}
