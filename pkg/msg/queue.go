// planner.go
package msg

import (
	"fmt"
	"os"
)

type Queue struct {
	Capacity int
	msg      chan interface{} // the channel of entries in planner
}

func NewQueue(capacity int) *Queue {

	chs := make(chan interface{}, capacity)

	return &Queue{
		Capacity: capacity,
		msg:      chs,
	}
}

// Add  entry in asynchronously
func (p *Queue) Send(entry interface{}) {
	fmt.Fprintf(os.Stderr, "Before adding entry %s to %s", entry, p.msg)
	p.msg <- entry
	fmt.Fprintf(os.Stderr, "After adding entry %s to %s", entry, p.msg)
}

// Recieve entry in asynchronously
func (p *Queue) Recieve() interface{} {
	fmt.Fprintf(os.Stderr, "waiting for entry in queue %s", p.msg)
	entry := <-p.msg
	fmt.Fprintf(os.Stderr, "Getting entry %s from %s", entry, p.msg)
	return entry
}
