package renderer

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/fmpwizard/go-quilljs-delta/delta"
)

const SampleSource = `package foo

import "fmt"

func PrintBar(a *int) error {
	fmt.Printf("Bar: %d\n", *a)
	return nil
}
`

const AlgolRenderedSource = `{"ops":[{"retain":7,"attributes":{"background":"#ffffff","bold":true,"underline":true}},{"retain":6,"attributes":{"background":"#ffffff"}},{"retain":6,"attributes":{"background":"#ffffff","bold":true,"underline":true}},{"retain":1,"attributes":{"background":"#ffffff"}},{"retain":5,"attributes":{"background":"#ffffff","color":"#666666","italic":true}},{"retain":2,"attributes":{"background":"#ffffff"}},{"retain":4,"attributes":{"background":"#ffffff","bold":true,"italic":true,"underline":true}},{"retain":1,"attributes":{"background":"#ffffff"}},{"retain":8,"attributes":{"background":"#ffffff","bold":true,"color":"#666666","italic":true}},{"retain":4,"attributes":{"background":"#ffffff"}},{"retain":3,"attributes":{"background":"#ffffff","bold":true,"underline":true}},{"retain":2,"attributes":{"background":"#ffffff"}},{"retain":5,"attributes":{"background":"#ffffff","bold":true,"underline":true}},{"retain":8,"attributes":{"background":"#ffffff"}},{"retain":6,"attributes":{"background":"#ffffff","bold":true,"color":"#666666","italic":true}},{"retain":1,"attributes":{"background":"#ffffff"}},{"retain":11,"attributes":{"background":"#ffffff","color":"#666666","italic":true}},{"retain":7,"attributes":{"background":"#ffffff"}},{"retain":6,"attributes":{"background":"#ffffff","bold":true,"underline":true}},{"retain":1,"attributes":{"background":"#ffffff"}},{"retain":3,"attributes":{"background":"#ffffff","bold":true,"underline":true}},{"retain":3,"attributes":{"background":"#ffffff"}}]}`

func debugDeltaString(t *testing.T, d delta.Delta) string {
	data, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestHighlighting(t *testing.T) {
	lexer := chroma.Coalesce(lexers.Get("go"))
	style := styles.Get("algol")
	iterator, err := lexer.Tokenise(nil, SampleSource)
	if err != nil {
		t.Fatal(err)
	}
	d, err := FormatToDelta(style, iterator)
	if err != nil {
		t.Fatal(err)
	}

	var d2 delta.Delta
	err = json.Unmarshal([]byte(AlgolRenderedSource), &d2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(d, d2) {
		t.Fatalf("Invalid change, expected: %s, actual: %s",
			debugDeltaString(t, d2), debugDeltaString(t, d))
	}
}

func TestChromaHighlighting(t *testing.T) {
	lexer := chroma.Coalesce(lexers.Get("go"))
	style := styles.Get("algol")
	iterator, err := lexer.Tokenise(nil, SampleSource)
	if err != nil {
		t.Fatal(err)
	}
	var buffer bytes.Buffer
	err = FormatToChroma(&buffer, style, iterator)
	if err != nil {
		t.Fatal(err)
	}

	var d delta.Delta
	err = json.Unmarshal(buffer.Bytes(), &d)
	if err != nil {
		t.Fatal(err)
	}
	var d2 delta.Delta
	err = json.Unmarshal([]byte(AlgolRenderedSource), &d2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(d, d2) {
		t.Fatalf("Invalid change, expected: %s, actual: %s",
			debugDeltaString(t, d2), debugDeltaString(t, d))
	}
}

func TestNullHighlighting(t *testing.T) {
	lexer := chroma.Coalesce(lexers.Get("go"))
	builder := chroma.NewStyleBuilder("null")
	style, _ := builder.Build()
	iterator, err := lexer.Tokenise(nil, SampleSource)
	if err != nil {
		t.Fatal(err)
	}
	d, err := FormatToDelta(style, iterator)
	if err != nil {
		t.Fatal(err)
	}

	d2 := *delta.New(nil).Retain(100, nil)
	if !reflect.DeepEqual(d, d2) {
		t.Fatalf("Invalid change, expected: %s, actual: %s",
			debugDeltaString(t, d2), debugDeltaString(t, d))
	}
}
