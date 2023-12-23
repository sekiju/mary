package utils

import (
	"fmt"
	"strings"
)

type IndexNameFormatter struct {
	pad string
}

func NewIndexNameFormatter(pagesCount int) *IndexNameFormatter {
	count := 999
	if pagesCount > 0 {
		count = pagesCount
	}
	pad := strings.Repeat("0", len(fmt.Sprint(count)))
	return &IndexNameFormatter{pad: pad}
}

func (pn *IndexNameFormatter) GetName(i int, ext string) string {
	str := fmt.Sprint(i)
	return pn.pad[:len(pn.pad)-len(str)] + str + ext
}
