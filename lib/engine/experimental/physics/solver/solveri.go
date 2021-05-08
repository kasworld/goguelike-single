package solver

import "github.com/kasworld/goguelike-single/lib/engine/experimental/physics/equation"

// SolverI is the interface type for all constraint solvers.
type SolverI interface {
	Solve(dt float32, nBodies int) *Solution // Solve should return the number of iterations performed
	AddEquation(eq equation.EquationI)
	RemoveEquation(eq equation.EquationI) bool
	ClearEquations()
}
