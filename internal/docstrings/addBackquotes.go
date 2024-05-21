package docstrings

import (
	"fmt"
	"strings"
)

func addBackquotes(s string) string {
	count := strings.Count(s, "%s")
	fmtArgs := make([]any, count)
	for i := 0; i < count; i++ {
		fmtArgs[i] = "`"
	}
	return fmt.Sprintf(s, fmtArgs...)
}
