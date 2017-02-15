package feedfinder

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestCheckLink(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
	}))
	defer ts.Close()

	tests := []struct {
		want *Link
	}{
		{
			want: &Link{
				URL:          ts.URL,
				ContentType:  "image/png",
				ResponseCode: 200,
				ETag:         "",
			},
		},
	}

	for _, test := range tests {
		client := New("testingnomatch")

		ch := make(chan *Link, 1)
		client.checkLink(ts.URL, ch)

		got := <-ch

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("want Link = %+v, got %+v", test.want, got)
		}
	}
}

func TestLinksFromURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
	}))
	defer ts.Close()

	tests := []struct {
		want []*Link
	}{
		{
			want: []*Link{
				{
					URL:          ts.URL,
					ContentType:  "image/png",
					ResponseCode: 200,
					ETag:         "",
				},
			},
		},
	}

	for _, test := range tests {
		client := New("testingnomatch")

		got := client.linksFromURL([]string{ts.URL})

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("want []Link = %+v, got %+v", test.want, got)
		}
	}
}

func TestDocumentFromReader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
	}))
	defer ts.Close()

	tests := []struct {
		want *Document
	}{
		{
			want: &Document{
				Links: []*Link{
					{
						URL:          ts.URL,
						ContentType:  "image/png",
						ResponseCode: 200,
						ETag:         "",
					},
				},
			},
		},
	}

	for _, test := range tests {
		client := &Client{
			Client: http.DefaultClient,
			DiscoverFunc: func(io.Reader) ([]string, error) {
				return []string{ts.URL}, nil
			},
			FilterFunc: func(s string, r []string, re *regexp.Regexp) ([]string, error) {
				return r, nil
			},
		}

		got, err := client.DocumentFromURL(ts.URL)
		if err != nil {
			t.Errorf("want err = nil, got %v", err)
			continue
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("want []Link = %+v, got %+v", test.want, got)
		}
	}
}

func TestDocumentFromURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
	}))
	defer ts.Close()

	tests := []struct {
		want *Document
	}{
		{
			want: &Document{
				Links: []*Link{
					{
						URL:          ts.URL,
						ContentType:  "image/png",
						ResponseCode: 200,
						ETag:         "",
					},
				},
			},
		},
	}

	for _, test := range tests {
		client := &Client{
			Client: http.DefaultClient,
			DiscoverFunc: func(io.Reader) ([]string, error) {
				return []string{ts.URL}, nil
			},
			FilterFunc: func(s string, r []string, re *regexp.Regexp) ([]string, error) {
				return r, nil
			},
		}

		got, err := client.DocumentFromReader(ts.URL, strings.NewReader(""))
		if err != nil {
			t.Errorf("want err = nil, got %+v", err)
			continue
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("want []Link = %+v, got %+v", test.want, got)
		}
	}
}
