package clients

import "strings"

type RequestOptions struct {
	Headers         map[string]string
	QueryParameters map[string]string
}

func DefaultRequestOptions() RequestOptions {
	return RequestOptions{
		Headers:         make(map[string]string),
		QueryParameters: make(map[string]string),
	}
}

func NewRequestOptions(headers map[string]string, queryParameters map[string][]string) RequestOptions {
	opts := DefaultRequestOptions()

	opts.Headers = headers

	for key, values := range queryParameters {
		opts.QueryParameters[key] = strings.Join(values, ",")
	}

	return opts
}
