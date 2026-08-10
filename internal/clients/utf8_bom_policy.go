package clients

import (
	"bufio"
	"bytes"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

var utf8BOM = []byte{0xef, 0xbb, 0xbf}

type utf8BOMPolicy struct{}

func newUTF8BOMPolicy() policy.Policy {
	return &utf8BOMPolicy{}
}

func (*utf8BOMPolicy) Do(req *policy.Request) (*http.Response, error) {
	resp, err := req.Next()
	if err != nil || resp.Body == nil || !isJSONContentType(resp.Header.Get("Content-Type")) {
		return resp, err
	}

	reader := bufio.NewReader(resp.Body)
	prefix, peekErr := reader.Peek(len(utf8BOM))
	if peekErr == nil && bytes.Equal(prefix, utf8BOM) {
		_, _ = reader.Discard(len(utf8BOM))
		if resp.ContentLength >= int64(len(utf8BOM)) {
			resp.ContentLength -= int64(len(utf8BOM))
		}
	}
	resp.Body = struct {
		io.Reader
		io.Closer
	}{
		Reader: reader,
		Closer: resp.Body,
	}

	return resp, nil
}

func isJSONContentType(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	mediaType = strings.ToLower(mediaType)
	return mediaType == "application/json" || strings.HasSuffix(mediaType, "+json")
}
