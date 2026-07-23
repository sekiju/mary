package util

import (
	"github.com/knst0/mdl/sdk/manga"
	"resty.dev/v3"
	"testing"
)

// AssertImage downloads page.URL, decodes it if page.Decode is set, and
// compares the result against a golden file committed under testdata/.
//
// Run tests with -update to regenerate the golden file from the current
// live response (review the diff before committing).
func AssertImage(t *testing.T, goldenPath string, page *manga.Page) {
	res, err := resty.New().R().Get(page.URL)
	if err != nil {
		t.Fatal(err)
		return
	}

	actual := res.Bytes()

	if page.Decode != nil {
		actual, err = page.Decode(actual)
		if err != nil {
			t.Fatal(err)
			return
		}
	}

	AssertGolden(t, goldenPath, actual)
}
