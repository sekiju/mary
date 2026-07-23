// Package testutil provides shared e2e snapshot-testing helpers for
// extractors. Tests hit real sites and hash downloaded page bytes instead of
// storing images, so fixtures stay small and don't redistribute copyrighted
// content.
package testutil

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/knst0/mdl/sdk/manga"
)

var Update = flag.Bool("update", false, "update snapshot fixtures instead of comparing against them")

type PageSnapshot struct {
	Filename string `json:"filename"`
	SHA256   string `json:"sha256"`
}

// HashPages downloads and decodes every page, returning one PageSnapshot per
// page in order.
func HashPages(t *testing.T, pages []*manga.Page) []PageSnapshot {
	t.Helper()

	client := &http.Client{}
	snapshots := make([]PageSnapshot, len(pages))

	for i, page := range pages {
		req, err := http.NewRequest(http.MethodGet, page.URL, nil)
		if err != nil {
			t.Fatalf("page %d: build request: %v", i, err)
		}
		for k, v := range page.Headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("page %d: fetch %s: %v", i, page.URL, err)
		}
		body, err := readAll(resp)
		if err != nil {
			t.Fatalf("page %d: read body: %v", i, err)
		}

		if page.Decode != nil {
			body, err = page.Decode(body)
			if err != nil {
				t.Fatalf("page %d: decode: %v", i, err)
			}
		}

		sum := sha256.Sum256(body)
		snapshots[i] = PageSnapshot{
			Filename: page.Filename,
			SHA256:   hex.EncodeToString(sum[:]),
		}
	}

	return snapshots
}

func readAll(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	buf := make([]byte, 0, 64*1024)
	for {
		chunk := make([]byte, 32*1024)
		n, err := resp.Body.Read(chunk)
		if n > 0 {
			buf = append(buf, chunk[:n]...)
		}
		if err != nil {
			if err.Error() == "EOF" {
				return buf, nil
			}
			return buf, err
		}
	}
}

// CompareSnapshot compares got against the JSON fixture at path, or writes
// got to path when the -update flag is set.
func CompareSnapshot(t *testing.T, path string, got []PageSnapshot) {
	t.Helper()

	if *Update {
		data, err := json.MarshalIndent(got, "", "  ")
		if err != nil {
			t.Fatalf("marshal snapshot: %v", err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatalf("write snapshot %s: %v", path, err)
		}
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read snapshot %s (run with -update to create it): %v", path, err)
	}

	var want []PageSnapshot
	if err := json.Unmarshal(data, &want); err != nil {
		t.Fatalf("parse snapshot %s: %v", path, err)
	}

	if len(want) != len(got) {
		t.Fatalf("page count mismatch: snapshot has %d, got %d", len(want), len(got))
	}

	for i := range want {
		if want[i].SHA256 != got[i].SHA256 {
			t.Errorf("page %d (%s): hash mismatch: snapshot=%s got=%s", i, got[i].Filename, want[i].SHA256, got[i].SHA256)
		}
		if want[i].Filename != got[i].Filename {
			t.Errorf("page %d: filename mismatch: snapshot=%s got=%s", i, want[i].Filename, got[i].Filename)
		}
	}

	if t.Failed() {
		fmt.Println("re-run with -update if this change is expected")
	}
}
