package shape

import "github.com/kasworld/goguelike-single/lib/engine/math32"

// ShapeI is the interface for all collision shapes.
// Shapes in this package satisfy this interface and also geometry.Geometry.
type ShapeI interface {
	BoundingBox() math32.Box3
	BoundingSphere() math32.Sphere
	Area() float32
	Volume() float32
	RotationalInertia(mass float32) math32.Matrix3
	ProjectOntoAxis(localAxis *math32.Vector3) (float32, float32)
}
