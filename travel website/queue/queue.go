.
package queue

import (
)

type CacheConfig EventQueue

type EventQueue interface {
	String() string
	
}
