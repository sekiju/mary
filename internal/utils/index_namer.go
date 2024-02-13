package utils

import (
	"fmt"
	"strings"
)

type IndexNamer struct {
	pad string
}

func NewIndexNamer(pagesCount int) *IndexNamer {
	count := 999
	if pagesCount > 0 {
		count = pagesCount
	}
	pad := strings.Repeat("0", len(fmt.Sprint(count)))
	return &IndexNamer{pad: pad}
}

func (pn *IndexNamer) Get(i int, ext string) string {
	str := fmt.Sprint(i)
	return pn.pad[:len(pn.pad)-len(str)] + str + ext
}
