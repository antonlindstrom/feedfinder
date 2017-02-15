package feedfinder

import (
	"io"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Discover finds the Atom or RSS feed from a HTML document.
func Discover(r io.Reader) ([]string, error) {
	var feeds []string

	z := html.NewTokenizer(r)

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if err := z.Err(); err == io.EOF {
				break
			}
			return feeds, nil
		}

		t := z.Token()
		if t.DataAtom == atom.Link && (t.Type == html.StartTagToken || t.Type == html.SelfClosingTagToken) {
			attrs := make(map[string]string)
			for _, a := range t.Attr {
				attrs[a.Key] = a.Val
			}

			if attrs["rel"] == "alternate" && attrs["href"] != "" && (attrs["type"] == "application/rss+xml" || attrs["type"] == "application/atom+xml") {
				feeds = append(feeds, attrs["href"])
			}
		}
	}

	return feeds, nil
}
