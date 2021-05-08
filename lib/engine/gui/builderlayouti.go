package gui

// BuilderLayoutI is the interface for all layout builders
type BuilderLayoutI interface {
	BuildLayout(b *Builder, am map[string]interface{}) (LayoutI, error)
	BuildParams(b *Builder, am map[string]interface{}) (interface{}, error)
}
