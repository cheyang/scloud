// planner.go
package msg

import (
	"fmt"
	"os"
	//	"sync"
)

type Queue struct {
	Capacity int
	msg      chan interface{} // the channel of entries in planner
	done     bool
	//	lock     *sync.Mutex
}

func NewQueue(capacity int) *Queue {

	chs := make(chan interface{}, capacity)

	return &Queue{
		Capacity: capacity,
		msg:      chs,
		done:     false,
	}
}

func (p *Queue) Length() int {
	return p.Capacity
}

// Add  entry in asynchronously
func (p *Queue) Send(entry interface{}) {
	fmt.Fprintf(os.Stderr, "Before adding entry %s to %v\n", entry, p.msg)
	p.msg <- entry
	fmt.Fprintf(os.Stderr, "After adding entry %s to %v\n", entry, p.msg)
}

// Recieve entry in asynchronously
func (p *Queue) Recieve() interface{} {
	fmt.Fprintf(os.Stderr, "waiting for entry in queue %v\n", p.msg)
	entry := <-p.msg
	fmt.Fprintf(os.Stderr, "Getting entry %s from %v\n", entry, p.msg)
	return entry
}

// Sender set done
func (p *Queue) SetDone() {
	fmt.Fprintf(os.Stderr, "Tell receiver it's done\n")
	p.done = true
}

func (p *Queue) IsDone() bool {
	return p.done
}

func (p *Queue) Close() {
	close(p.msg)
}
