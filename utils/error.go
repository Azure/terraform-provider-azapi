package utils

import (
	"fmt"
	"strings"
)

func ErrorMismatch(key, expected, actual string) error {
	if strings.HasPrefix(key, ".") {
		key = key[1:]
	}
	return fmt.Errorf("`%s` is invalid, expect `%s` but got `%s`", key, expected, actual)
}

func ErrorNotMatchAny(key string) error {
	if strings.HasPrefix(key, ".") {
		key = key[1:]
	}
	return fmt.Errorf("`%s` doesn't match any accepted values", key)
}

func ErrorShouldNotDefine(key string, options []string) error {
	if strings.HasPrefix(key, ".") {
		key = key[1:]
	}
	suggestion := getSuggestion(key, options)
	if strings.HasPrefix(suggestion, ".") {
		suggestion = suggestion[1:]
	}
	return fmt.Errorf("`%s` is not expected here. Do you mean `%s`? ", key, suggestion)
}

func ErrorShouldDefine(key string) error {
	if strings.HasPrefix(key, ".") {
		key = key[1:]
	}
	return fmt.Errorf("`%s` is required, but no definition was found", key)
}

func getSuggestion(value string, options []string) string {
	suggestion := ""
	distance := 1 << 16
	for _, option := range options {
		if dist := editDistance(value, option); dist < distance {
			distance = dist
			suggestion = option
		}
	}
	return suggestion
}

func editDistance(a, b string) int {
	n, m := len(a), len(b)
	f := make([][]int, n+1)
	for i := range f {
		f[i] = make([]int, m+1, 1<<16)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			f[i][j] = 1 << 16
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				f[i][j] = f[i-1][j-1]
			}
			if f[i][j] > f[i-1][j]+1 {
				f[i][j] = f[i-1][j] + 1
			}
			if f[i][j] > f[i][j-1]+1 {
				f[i][j] = f[i][j-1] + 1
			}
			if f[i][j] > f[i-1][j-1]+1 {
				f[i][j] = f[i-1][j-1] + 1
			}
		}
	}
	return f[n][m]
}
