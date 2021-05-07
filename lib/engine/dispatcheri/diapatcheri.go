package dispatcheri

type EventName string

// Callback is the type for Dispatcher callback functions.
type Callback func(evname EventName, ev interface{})

// DispatcherI is the interface for event dispatchers.
type DispatcherI interface {
	Subscribe(evname EventName, cb Callback)
	SubscribeID(evname EventName, id interface{}, cb Callback)
	UnsubscribeID(evname EventName, id interface{}) int
	UnsubscribeAllID(id interface{}) int
	Dispatch(evname EventName, ev interface{}) int
}
