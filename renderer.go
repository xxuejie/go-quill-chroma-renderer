package renderer

import (
	"encoding/json"
	"fmt"
	"io"
	"math"

	"github.com/alecthomas/chroma"
	"github.com/fmpwizard/go-quilljs-delta/delta"
)

func FormatToDelta(style *chroma.Style, it chroma.Iterator) (delta.Delta, error) {
	d := delta.New(nil)
	for token := it(); token != chroma.EOF; token = it() {
		attributes := make(map[string]interface{})
		entry := style.Get(token.Type)
		if !entry.IsZero() {
			if entry.Bold == chroma.Yes {
				attributes["bold"] = true
			}
			if entry.Underline == chroma.Yes {
				attributes["underline"] = true
			}
			if entry.Italic == chroma.Yes {
				attributes["italic"] = true
			}
			if entry.Colour.IsSet() {
				attributes["color"] = fmt.Sprintf("#%02x%02x%02x", entry.Colour.Red(),
					entry.Colour.Green(), entry.Colour.Blue())
			}
			if entry.Background.IsSet() {
				attributes["background"] = fmt.Sprintf("#%02x%02x%02x",
					entry.Background.Red(),
					entry.Background.Green(),
					entry.Background.Blue())
			}
		}
		if len(attributes) == 0 {
			attributes = nil
		}
		d = d.Retain(len([]rune(token.Value)), attributes)
	}
	return *d, nil
}

func FormatToChroma(w io.Writer, style *chroma.Style, it chroma.Iterator) error {
	delta, err := FormatToDelta(style, it)
	if err != nil {
		return err
	}
	data, err := json.Marshal(delta)
	if err != nil {
		return err
	}
	if len(data) > math.MaxUint32 {
		return fmt.Errorf("Delta too long: %d", len(data))
	}
	length := fmt.Sprintf("%11d ", len(data))
	_, err = w.Write([]byte(length))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
