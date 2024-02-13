package utils

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

func JoinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

func ExportURLQueries(path string) (string, url.Values) {
	queries := make(map[string][]string)

	parts := strings.Split(path, "?")
	if len(parts) == 2 {
		queryPart := parts[1]
		parsedQuery, err := url.ParseQuery(queryPart)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return "", nil
		}

		for key, values := range parsedQuery {
			if len(values) > 0 {
				queries[key] = values
			}
		}
	}

	return parts[0], queries
}

func LastURLSegment(path string) string {
	path = strings.TrimSuffix(path, "/")
	return filepath.Base(path)
}
