package equation

import "github.com/kasworld/goguelike-single/lib/engine/math32"

// BodyI is the interface of all body types.
type BodyI interface {
	Index() int
	Position() math32.Vector3
	Velocity() math32.Vector3
	AngularVelocity() math32.Vector3
	Force() math32.Vector3
	Torque() math32.Vector3
	InvMassEff() float32
	InvRotInertiaWorldEff() *math32.Matrix3
}

// EquationI is the interface type for all equations types.
type EquationI interface {
	SetBodyA(BodyI)
	BodyA() BodyI
	SetBodyB(BodyI)
	BodyB() BodyI
	JeA() JacobianElement
	JeB() JacobianElement
	SetEnabled(state bool)
	Enabled() bool
	MinForce() float32
	MaxForce() float32
	Eps() float32
	SetMultiplier(multiplier float32)
	ComputeB(h float32) float32
	ComputeC() float32
}
