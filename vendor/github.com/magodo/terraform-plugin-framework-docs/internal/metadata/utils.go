package metadata

import (
	"context"
	"strings"
	"unicode"
)

type DescriptionProvider interface {
	GetDescription() string
}

type MarkdownDescriptionProvider interface {
	GetMarkdownDescription() string
}

type BothDescriptionProvider interface {
	DescriptionProvider
	MarkdownDescriptionProvider
}

func DescriptionOf[T BothDescriptionProvider](d T) string {
	if v, ok := any(d).(MarkdownDescriptionProvider); ok && v.GetMarkdownDescription() != "" {
		return v.GetMarkdownDescription()
	}
	return any(d).(DescriptionProvider).GetDescription()
}

func MaybeDescriptionOf(d any) *string {
	if v, ok := d.(MarkdownDescriptionProvider); ok && v.GetMarkdownDescription() != "" {
		return new(v.GetMarkdownDescription())
	}
	if v, ok := d.(DescriptionProvider); ok && v.GetDescription() != "" {
		return new(v.GetDescription())
	}
	return nil
}

type DescriptionCtxProvider interface {
	Description(context.Context) string
}

type MarkdownDescriptionCtxProvider interface {
	MarkdownDescription(context.Context) string
}

type BothDescriptionCtxProvider interface {
	DescriptionCtxProvider
	MarkdownDescriptionCtxProvider
}

func DescriptionCtxOf[T BothDescriptionCtxProvider](ctx context.Context, d T) string {
	if v, ok := any(d).(MarkdownDescriptionCtxProvider); ok && v.MarkdownDescription(ctx) != "" {
		return v.MarkdownDescription(ctx)
	}
	return any(d).(DescriptionCtxProvider).Description(ctx)
}

func MaybeDescriptionCtxOf(ctx context.Context, d any) *string {
	if v, ok := d.(MarkdownDescriptionCtxProvider); ok && v.MarkdownDescription(ctx) != "" {
		return new(v.MarkdownDescription(ctx))
	}
	if v, ok := d.(DescriptionCtxProvider); ok && v.Description(ctx) != "" {
		return new(v.Description(ctx))
	}
	return nil
}

func MapSlice[T any, U any](input []T, f func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = f(v)
	}
	return result
}

func MapSliceSome[T any, U any](input []T, f func(T) *U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		ptr := f(v)
		if ptr != nil {
			result[i] = *ptr
		}
	}
	return result
}

func MapOrNil[T any, U any](input T, f func(T) U) *U {
	var anyInput any = input
	if anyInput == nil {
		return nil
	}
	output := f(input)
	return &output
}

func MapOrZero[T any, U any](input T, f func(T) U) U {
	var zero U
	var anyInput any = input
	if anyInput == nil {
		return zero
	}
	output := f(input)
	return output
}

func PointerTo[T any](ptr *T) T {
	var zero T
	if ptr == nil {
		return zero
	}
	return *ptr
}

func Sentencefy(s string) string {
	if s == "" {
		return s
	}
	return capitalizeFirstLetter(ensureStringEndsWithDot(s))
}

func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}

	// Convert string to a slice of runes to handle Unicode characters correctly
	runes := []rune(s)

	// Capitalize the first rune
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func ensureStringEndsWithDot(s string) string {
	if s == "" {
		return s
	}
	return strings.TrimRight(s, ".") + "."
}
