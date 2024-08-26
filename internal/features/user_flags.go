package features

type UserFeatures struct {
	DefaultTags     map[string]string
	DefaultLocation string
	DefaultNaming   string
}

func Default() UserFeatures {
	return UserFeatures{
		DefaultTags:     nil,
		DefaultLocation: "",
		DefaultNaming:   "",
	}
}
