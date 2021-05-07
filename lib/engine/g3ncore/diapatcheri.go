package g3ncore

// DispatcherI is the interface for event dispatchers.
type DispatcherI interface {
	Subscribe(evname string, cb Callback)
	SubscribeID(evname string, id interface{}, cb Callback)
	UnsubscribeID(evname string, id interface{}) int
	UnsubscribeAllID(id interface{}) int
	Dispatch(evname string, ev interface{}) int
}
