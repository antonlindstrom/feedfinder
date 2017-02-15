package feedfinder

import (
	"regexp"
)

// DefaultFilter returns the resources that doesn't match the regexp (it
// filters them out) and also adds absolute URLs to relative ones.
func DefaultFilter(url string, r []string, re *regexp.Regexp) ([]string, error) {
	var resources []string

	for _, res := range r {
		res, err := ToAbs(url, res)
		if err != nil {
			return nil, err
		}

		if re.MatchString(res) {
			continue
		}

		resources = append(resources, res)
	}

	return resources, nil
}
