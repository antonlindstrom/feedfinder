package feedfinder

import (
	"net/url"
	"path"
	"strings"
)

// ToAbs returns the absolute value for a resource.
func ToAbs(fullURL, resource string) (string, error) {
	if isAbs(resource) {
		return resource, nil
	}
	return join(baseURL(fullURL), resource)
}

// isAbs checks if it's a resource starts with http.
func isAbs(resource string) bool {
	return strings.HasPrefix(resource, "http")
}

// join joins the base url and a resource path.
func join(baseURL, resource string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, resource)

	return u.String(), nil
}

// baseURL returns the absolute URL for a link, if the url.Parse fails,
// returns an empty string.
func baseURL(fullURL string) string {
	u, err := url.Parse(fullURL)
	if err != nil || (u.Scheme == "" && u.Host == "") {
		return ""
	}

	return u.Scheme + "://" + u.Host
}
