package feedfinder

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

// Client contains an HTTP client and a URL parser.
type Client struct {
	*http.Client
	regexpFilter *regexp.Regexp
	DiscoverFunc func(io.Reader) ([]string, error)
	FilterFunc   func(string, []string, *regexp.Regexp) ([]string, error)
}

// Document contains data about the document that should be checked.
type Document struct {
	Links []*Link
}

// Link contains status about a link.
type Link struct {
	URL          string
	ContentType  string
	ETag         string
	ResponseCode int
	Error        error
}

// New creates a new default client and parser.
func New(regexpFilter string) *Client {
	return &Client{
		Client:       http.DefaultClient,
		regexpFilter: regexp.MustCompile(regexpFilter),
		DiscoverFunc: Discover,
		FilterFunc:   DefaultFilter,
	}
}

// DocumentFromURL fetches the body of a URL and returns a document with the
// links set.
func (c *Client) DocumentFromURL(url string) (*Document, error) {
	if strings.HasPrefix(url, "//") {
		url = "https:" + url
	}

	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}

	doc, err := c.DocumentFromReader(url, resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, resp.Body.Close()
}

// DocumentFromReader populates a document from an io.Reader.
func (c *Client) DocumentFromReader(u string, r io.Reader) (*Document, error) {
	resources, err := c.DiscoverFunc(r)
	if err != nil {
		return nil, err
	}

	resources, err = c.FilterFunc(u, resources, c.regexpFilter)
	if err != nil {
		return nil, err
	}

	return &Document{Links: c.linksFromURL(resources)}, nil
}

// linksFromURL fetches all the links and populates the Link struct.
// This does a concurrent fetch of all the urls.
func (c *Client) linksFromURL(urls []string) []*Link {
	ch := make(chan *Link, len(urls))

	for _, url := range urls {
		go c.checkLink(url, ch)
	}

	var links []*Link
	for i := 0; i < len(urls); i++ {
		links = append(links, <-ch)
	}

	close(ch)
	return links
}

// checkLink checks the status of a link, which translates to doing a HEAD
// request to the url.
func (c *Client) checkLink(url string, ch chan<- *Link) {
	resp, err := c.Head(url)
	if err != nil {
		ch <- &Link{
			URL:   url,
			Error: err,
		}
		return
	}

	ch <- &Link{
		URL:          url,
		ContentType:  resp.Header.Get("Content-Type"),
		ETag:         resp.Header.Get("ETag"),
		ResponseCode: resp.StatusCode,
	}
}
