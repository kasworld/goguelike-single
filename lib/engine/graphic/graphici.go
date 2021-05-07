package graphic

import (
	"github.com/kasworld/goguelike-single/lib/engine/g3ncore"
	"github.com/kasworld/goguelike-single/lib/engine/geometry"
	"github.com/kasworld/goguelike-single/lib/engine/gls"
	"github.com/kasworld/goguelike-single/lib/engine/renderinfo"
)

// GraphicI is the interface for all Graphic objects.
type GraphicI interface {
	g3ncore.NodeI
	GetGraphic() *Graphic
	GetGeometry() *geometry.Geometry
	GeometryI() geometry.GeometryI
	SetRenderable(bool)
	Renderable() bool
	SetCullable(bool)
	Cullable() bool
	RenderSetup(gs *gls.GLS, rinfo *renderinfo.RenderInfo)
}
