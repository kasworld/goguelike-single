// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dispatcher

import (
	"github.com/kasworld/goguelike-single/lib/engine/dispatcheri"
	"github.com/kasworld/goguelike-single/lib/engine/eventtype"
)

// Dispatcher implements an event dispatcher.
type Dispatcher struct {
	evmap map[eventtype.EventType][]subscription // Map of event names to subscription lists
}

// subscription links a dispatcheri.Callback with a user-provided unique id.
type subscription struct {
	id       interface{}
	callBack dispatcheri.Callback
}

// NewDispatcher creates and returns a new event dispatcher.
func NewDispatcher() *Dispatcher {

	d := new(Dispatcher)
	d.Initialize()
	return d
}

// Initialize initializes the event dispatcher.
// It is normally used by other types which embed a dispatcher.
func (d *Dispatcher) Initialize() {

	d.evmap = make(map[eventtype.EventType][]subscription)
}

// Subscribe subscribes a callback to events with the given name.
// If it is necessary to unsubscribe later, SubscribeID should be used instead.
func (d *Dispatcher) Subscribe(evname eventtype.EventType, callBack dispatcheri.Callback) {

	d.evmap[evname] = append(d.evmap[evname], subscription{nil, callBack})
}

// SubscribeID subscribes a callback to events events with the given name.
// The user-provided unique id can be used to unsubscribe via UnsubscribeID.
func (d *Dispatcher) SubscribeID(evname eventtype.EventType, id interface{}, callBack dispatcheri.Callback) {

	d.evmap[evname] = append(d.evmap[evname], subscription{id, callBack})
}

// UnsubscribeID removes all subscribed callbacks with the specified unique id
//	from the specified event.
// Returns the number of subscriptions removed.
func (d *Dispatcher) UnsubscribeID(evname eventtype.EventType, id interface{}) int {

	// Get list of subscribers for this event
	subs := d.evmap[evname]
	if len(subs) == 0 {
		return 0
	}

	// Remove all subscribers of the specified event with the specified id,
	// counting how many were removed
	rm := 0
	i := 0
	for _, s := range subs {
		if s.id == id {
			rm++
		} else {
			subs[i] = s
			i++
		}
	}
	d.evmap[evname] = subs[:i]
	return rm
}

// UnsubscribeAllID removes all subscribed callbacks with the specified unique id from all events.
// Returns the number of subscriptions removed.
func (d *Dispatcher) UnsubscribeAllID(id interface{}) int {

	// Remove all subscribers with the specified id (for all events), counting how many were removed
	total := 0
	for evname := range d.evmap {
		total += d.UnsubscribeID(evname, id)
	}
	return total
}

// Dispatch dispatches the specified event to all registered subscribers.
// The function returns the number of subscribers to which the event was dispatched.
func (d *Dispatcher) Dispatch(evname eventtype.EventType, ev interface{}) int {

	// Get list of subscribers for this event
	subs := d.evmap[evname]
	nsubs := len(subs)
	if nsubs == 0 {
		return 0
	}

	// Dispatch event to all subscribers
	for _, s := range subs {
		s.callBack(evname, ev)
	}
	return nsubs
}
