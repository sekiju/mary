package util

import (
	"fmt"
	json "github.com/bytedance/sonic"
	"strings"
)

func ExtractStringFromHTML(html, prefix, suffix string) (string, error) {
	startIdx := strings.Index(html, prefix)
	if startIdx == -1 {
		return "", fmt.Errorf("prefix not found: %s", prefix)
	}

	startIdx += len(prefix)

	endIdx := strings.Index(html[startIdx:], suffix)
	if endIdx == -1 {
		return "", fmt.Errorf("suffix not found: %s", suffix)
	}

	return html[startIdx : startIdx+endIdx], nil
}

func ExtractJSONFromHTML[T any](html, prefix, suffix string) (*T, error) {
	html = strings.ReplaceAll(html, "&quot;", "\"")

	escapedJSON, err := ExtractStringFromHTML(html, prefix, suffix)
	if err != nil {
		return nil, err
	}

	var result T
	if err = json.Unmarshal([]byte(escapedJSON), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
