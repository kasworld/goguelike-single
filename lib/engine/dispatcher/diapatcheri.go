package dispatcher

// Callback is the type for Dispatcher callback functions.
type Callback func(evname string, ev interface{})

// DispatcherI is the interface for event dispatchers.
type DispatcherI interface {
	Subscribe(evname string, cb Callback)
	SubscribeID(evname string, id interface{}, cb Callback)
	UnsubscribeID(evname string, id interface{}) int
	UnsubscribeAllID(id interface{}) int
	Dispatch(evname string, ev interface{}) int
}
