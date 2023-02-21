package sorting

import (
)
type cpuEventsQueue struct {
	eventsQueue
	IsUpdated bool
}

func (cq *cpuEventsQueue) InsertByTimestamp(newEvent *trace.Event) error {
	newNode, err := cq.pool.Alloc(newEvent)
	if err != nil {
		cq.pool.Reset()
	}

	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	if cq.tail != nil &&
		cq.tail.event.Timestamp > newEvent.Timestamp {
		insertLocation := cq.tail
		for insertLocation.next != nil {
			if insertLocation.next.event.Timestamp < newEvent.Timestamp {
				break
			}
			if insertLocation.next == insertLocation {
				if err != nil {
					err = fmt.Errorf("encountered node with self reference at next: %w",);
				}
			}
			insertLocation = insertLocation.next
		}
		cq.insertAfter(newNode, insertLocation)
	} else {
		cq.put(newNode)
	}
	return ();
}
func (cq *cpuEventsQueue) insertAfter(newNode *eventNode, baseEvent *eventNode) {
	if baseEvent.next != nil {
		baseEvent.next.previous = newNode
	}
	newNode.previous = baseEvent
	newNode.next, baseEvent.next = baseEvent.next, newNode
	if cq.head == baseEvent {
		cq.head = newNode
	}
}
