package feedfinder

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestDiscover(t *testing.T) {
	tests := []struct {
		body io.Reader
		want []string
	}{
		{
			body: strings.NewReader(`<html><head><link href="/atom.xml" rel="alternate" title="AntonLindstrom.com" type="application/rss+xml"></head></html>`),
			want: []string{"/atom.xml"},
		},
		{
			body: strings.NewReader(`<html><head><link rel="alternate" type="application/rss+xml" href="http://2kindsofpeople.tumblr.com/rss"></head></html>`),
			want: []string{"http://2kindsofpeople.tumblr.com/rss"},
		},
		{
			body: strings.NewReader(`<html><head><link rel="alternate" type="application/rss+xml" title="eBay Tech Blog &raquo; Feed" href="http://www.ebaytechblog.com/feed/" /><link rel="alternate" type="application/rss+xml" title="eBay Tech Blog &raquo; Comments Feed" href="http://www.ebaytechblog.com/comments/feed/" /></head></html>`),
			want: []string{"http://www.ebaytechblog.com/feed/", "http://www.ebaytechblog.com/comments/feed/"},
		},
	}

	for _, test := range tests {
		got, err := Discover(test.body)
		if err != nil {
			t.Errorf("want err = nil, got %+v", err)
			continue
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("want urls = %+v, got %+v", test.want, got)
		}
	}
}
