package metadata

type Categorizer interface {
	Category() Category
}

type Category string

const (
	CategoryProvider   = "Provider"
	CategoryResource   = "Resource"
	CategoryDataSource = "Data Source"
	CategoryEphemeral  = "Ephemeral Resource"
	CategoryAction     = "Action"
	CategoryList       = "List Resource"
	CategoryFunction   = "Function"
)
