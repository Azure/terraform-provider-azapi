package features

type UserFeatures struct {
	SchemaValidationEnabled bool
}

func Default() UserFeatures {
	return UserFeatures{
		SchemaValidationEnabled: true,
	}
}
