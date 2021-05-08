package dispatcheri

import "github.com/kasworld/goguelike-single/lib/engine/eventenum"

// Callback is the type for Dispatcher callback functions.
type Callback func(evname eventenum.EventName, ev interface{})

// DispatcherI is the interface for event dispatchers.
type DispatcherI interface {
	Subscribe(evname eventenum.EventName, cb Callback)
	SubscribeID(evname eventenum.EventName, id interface{}, cb Callback)
	UnsubscribeID(evname eventenum.EventName, id interface{}) int
	UnsubscribeAllID(id interface{}) int
	Dispatch(evname eventenum.EventName, ev interface{}) int
}
