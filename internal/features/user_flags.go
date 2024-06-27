package features

type UserFeatures struct {
	DefaultTags                  map[string]string
	DefaultLocation              string
	DefaultNaming                string
	DefaultNamingPrefix          string
	DefaultNamingSuffix          string
	CafEnabled                   bool
	EnableHCLOutputForDataSource bool
	EnablePreflight              bool
}

func Default() UserFeatures {
	return UserFeatures{
		DefaultTags:                  nil,
		DefaultLocation:              "",
		DefaultNaming:                "",
		DefaultNamingPrefix:          "",
		DefaultNamingSuffix:          "",
		CafEnabled:                   false,
		EnableHCLOutputForDataSource: false,
		EnablePreflight:              false,
	}
}
