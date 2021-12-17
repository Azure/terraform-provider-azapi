package features

type UserFeatures struct {
	SchemaValidationEnabled bool
	DefaultTags             map[string]string
}

func Default() UserFeatures {
	return UserFeatures{
		SchemaValidationEnabled: true,
		DefaultTags:             nil,
	}
}
