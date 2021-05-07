package geometry

import "github.com/kasworld/goguelike-single/lib/engine/gls"

// GeometryI is the interface for all geometries.
type GeometryI interface {
	GetGeometry() *Geometry
	RenderSetup(gs *gls.GLS)
	Dispose()
}
