package specifications

type Specification interface {
	ToSQL() (string, map[string]any)
}
