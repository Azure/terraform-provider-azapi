package location

import (
	"strings"
)

func Normalize(input string) string {
	return strings.ReplaceAll(strings.ToLower(input), " ", "")
}
