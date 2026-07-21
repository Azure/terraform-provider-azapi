package parse

import "testing"

func Test_hasIdentifierSegment(t *testing.T) {
	testData := []struct {
		Name         string
		ResourceType string
		Identifier   string
		Expected     bool
	}{
		{
			Name:         "name identifier present with api version suffix",
			ResourceType: "Microsoft.KeyVault/vaults/certificates@2016-10-01",
			Identifier:   "name",
			Expected:     true,
		},
		{
			Name:         "name identifier present without api version suffix",
			ResourceType: "Microsoft.KeyVault/vaults/certificates",
			Identifier:   "name",
			Expected:     true,
		},
		{
			Name:         "parentId identifier present",
			ResourceType: "Microsoft.KeyVault/vaults/certificates@2016-10-01",
			Identifier:   "parentId",
			Expected:     true,
		},
		{
			Name:         "identifier absent from url format",
			ResourceType: "Microsoft.KeyVault/vaults/certificates/contacts@v1.1-preview.2",
			Identifier:   "name",
			Expected:     false,
		},
		{
			Name:         "default form identifier ({name=default})",
			ResourceType: "Microsoft.Purview/accounts/Scanning/datasources/scans/triggers@2022-02-01-preview",
			Identifier:   "name",
			Expected:     true,
		},
		{
			Name:         "embedded identifier (indexes('{name}'))",
			ResourceType: "Microsoft.Search/searchServices/indexes@2023-11-01",
			Identifier:   "name",
			Expected:     true,
		},
		{
			Name:         "identifier prefix does not match longer placeholder",
			ResourceType: "Microsoft.KeyVault/vaults/certificates@2016-10-01",
			Identifier:   "nam",
			Expected:     false,
		},
		{
			Name:         "unmapped type falls back to true",
			ResourceType: "Microsoft.Foo/bar@2020-01-01",
			Identifier:   "name",
			Expected:     true,
		},
	}

	for _, v := range testData {
		t.Run(v.Name, func(t *testing.T) {
			actual := hasIdentifierSegment(v.ResourceType, v.Identifier)
			if actual != v.Expected {
				t.Fatalf("expected %v but got %v for resource type %q and identifier %q", v.Expected, actual, v.ResourceType, v.Identifier)
			}
		})
	}
}
