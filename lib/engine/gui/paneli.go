package gui

import "github.com/kasworld/goguelike-single/lib/engine/graphic"

// PanelI is the interface for all panel types
type PanelI interface {
	graphic.GraphicI
	GetPanel() *Panel
	Width() float32
	Height() float32
	Enabled() bool
	SetEnabled(bool)
	SetLayout(LayoutI)
	InsideBorders(x, y float32) bool
	SetZLayerDelta(zLayerDelta int)
	ZLayerDelta() int

	// TODO these methods here should probably be defined in NodeI
	SetPosition(x, y float32)
	SetPositionX(x float32)
	SetPositionY(y float32)
	SetPositionZ(y float32)
}
