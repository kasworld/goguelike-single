package material

import "github.com/kasworld/goguelike-single/lib/engine/gls"

// MaterialI is the interface for all materials.
type MaterialI interface {
	GetMaterial() *Material
	RenderSetup(gs *gls.GLS)
	Dispose()
}
