package features

type UserFeatures struct {
	DefaultTags                  map[string]string
	DefaultLocation              string
	DefaultNaming                string
	CafEnabled                   bool
	EnableHCLOutputForDataSource bool
}

func Default() UserFeatures {
	return UserFeatures{
		DefaultTags:                  nil,
		DefaultLocation:              "",
		DefaultNaming:                "",
		CafEnabled:                   false,
		EnableHCLOutputForDataSource: false,
	}
}
