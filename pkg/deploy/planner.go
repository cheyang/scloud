// planner.go
package deploy

import (
	"fmt"
	"os"
)

type Planner struct {
	Capacity  int
	recievers chan interface{} // the channel of entries in planner
}

func NewPlanner(capacity int) *Planner {

	chs := make(chan interface{}, capacity)

	return &Planner{
		Capacity:  capacity,
		recievers: chs,
	}
}

// Add  entry in asynchronously
func (p *Planer) Add(entry interface{}) {
	fmt.Fprintf(os.Stderr, "Before adding entry %s", entry)
	recievers <- entry
	fmt.Fprintf(os.Stderr, "After adding entry %s", entry)
}

// Recieve entry in asynchronously
func (p *Planer) Get() interface{} {
	fmt.Fprintf(os.Stderr, "Before adding entry %s", entry)
	entry := <-recievers
	fmt.Fprintf(os.Stderr, "After adding entry %s", entry)
	return entry
}
