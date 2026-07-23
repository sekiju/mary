package renamer

import (
	"fmt"
	"strings"
)

type Renamer struct {
	pad string
}

func New(fileCount int) *Renamer {
	var count int
	if fileCount > 0 {
		count = fileCount
	} else {
		count = 999
	}

	return &Renamer{pad: strings.Repeat("0", len(fmt.Sprint(count)))}
}

func (r *Renamer) Name(index int, extension string) string {
	str := fmt.Sprint(index)
	return r.pad[:len(r.pad)-len(str)] + str + extension
}
