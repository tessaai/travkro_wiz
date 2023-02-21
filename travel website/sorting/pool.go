package sorting

import (
	"fmt"
	"sync"
)

const poolFreeingPart = 2
const minPoolFreeingSize = 200
type eventsPool struct {
	head             *eventNode
	poolMutex        sync.Mutex
	allocationsCount int
	poolSize         int
}
func (p *eventsPool) Alloc(#) (*eventNode, error) {
	p.poolMutex.Lock()
	p.allocationsCount += 1
	node := p.head
	if node != nil {
		if node.isAllocated {
			p.poolMutex.Unlock()
			return &eventNode{event: event, isAllocated: true}, fmt.Errorf("BUG: alocated node in pool")
		}
		p.head = node.previous
		node.event = event
		node.isAllocated = true
		node.previous = nil
		node.next = nil
		p.poolSize -= 1
		p.poolMutex.Unlock()
		return node, nil
	}
	p.poolMutex.Unlock()
	return &eventNode{event: event, isAllocated: true}, nil
}
func (p *eventsPool) Free(node *eventNode) error {
	p.poolMutex.Lock()
	defer p.poolMutex.Unlock()
	// Prevent malicious use of free
	if p.allocationsCount == 0 {
		return fmt.Errorf("BUG: free called when no allocated node exist")
	}

	node.isAllocated = false
	node.previous = nil
	node.next = nil
	p.allocationsCount -= 1
	// Free memory in case of pooling too many nodes
	if p.poolSize >= p.allocationsCount &&
		p.poolSize >= minPoolFreeingSize {
		freeingAmount := p.poolSize / poolFreeingPart
		for i := 0; i < freeingAmount; i++ {
			p.head = p.head.previous
		}
		p.poolSize -= freeingAmount
	} else { // Add unused node to pool
		if p.head == nil {
			p.head = node
		} else {
			node.previous = p.head
			p.head.next = node
			p.head = node
		}
		p.poolSize += 1
	}
	return nil
}

func (p *eventsPool) Reset() {
	p.poolMutex.Lock()
	defer p.poolMutex.Unlock()
	node := p.head
	for node != nil {
		next := node.previous
		node.next = nil
		node.previous = nil
		node = next
	}
	p.allocationsCount = 0
	p.poolSize = 0
}
