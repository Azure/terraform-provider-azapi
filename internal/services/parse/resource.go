package parse

import (
	"fmt"
	"net/url"
	"strings"
)

type ResourceId struct {
	Url        string
	ApiVersion string
}

func NewResourceID(url, apiVersion string) ResourceId {
	return ResourceId{
		Url:        url,
		ApiVersion: apiVersion,
	}
}

func (id ResourceId) String() string {
	segments := []string{
		fmt.Sprintf("ResourceId %q", id.Url),
		fmt.Sprintf("Api Version %q", id.ApiVersion),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource", segmentsStr)
}

func (id ResourceId) ID() string {
	fmtString := "%s?api-version=%s"
	return fmt.Sprintf(fmtString, id.Url, id.ApiVersion)
}

// ResourceID parses a Resource ID into an ResourceId struct
func ResourceID(input string) (*ResourceId, error) {
	idUrl, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	resourceId := ResourceId{
		Url:        idUrl.Path,
		ApiVersion: idUrl.Query().Get("api-version"),
	}

	if resourceId.Url == "" {
		return nil, fmt.Errorf("ID was missing the 'url' element")
	}

	if resourceId.ApiVersion == "" {
		return nil, fmt.Errorf("ID was missing the 'api-version' element")
	}

	return &resourceId, nil
}
