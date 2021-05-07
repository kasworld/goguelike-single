package camerai

import "github.com/kasworld/goguelike-single/lib/engine/math32"

// CameraI is the interface for all cameras.
type CameraI interface {
	ViewMatrix(m *math32.Matrix4)
	ProjMatrix(m *math32.Matrix4)
}
