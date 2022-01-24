package features

type UserFeatures struct {
	SchemaValidationEnabled bool
	DefaultTags             map[string]string
	DefaultLocation         string
}

func Default() UserFeatures {
	return UserFeatures{
		SchemaValidationEnabled: true,
		DefaultTags:             nil,
		DefaultLocation:         "",
	}
}
