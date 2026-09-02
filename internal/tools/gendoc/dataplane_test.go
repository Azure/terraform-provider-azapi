package main

import "testing"

func TestEscapeTemplateDelimiters(t *testing.T) {
	const content = `content = "{{item.input}}"`
	const expected = `content = "{{"{{"}}item.input}}"`

	if actual := escapeTemplateDelimiters(content); actual != expected {
		t.Fatalf("unexpected escaped content: got %q, want %q", actual, expected)
	}
}
