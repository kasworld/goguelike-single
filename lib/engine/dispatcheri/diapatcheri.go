package dispatcheri

import "github.com/kasworld/goguelike-single/lib/engine/eventtype"

// Callback is the type for Dispatcher callback functions.
type Callback func(evname eventtype.EventType, ev interface{})

// DispatcherI is the interface for event dispatchers.
type DispatcherI interface {
	Subscribe(evname eventtype.EventType, cb Callback)
	SubscribeID(evname eventtype.EventType, id interface{}, cb Callback)
	UnsubscribeID(evname eventtype.EventType, id interface{}) int
	UnsubscribeAllID(id interface{}) int
	Dispatch(evname eventtype.EventType, ev interface{}) int
}
