package features

type UserFeatures struct {
	DefaultTags     map[string]string
	DefaultLocation string
}

func Default() UserFeatures {
	return UserFeatures{
		DefaultTags:     nil,
		DefaultLocation: "",
	}
}
