// State.go
package state

type State int

const (
	None = iota
	Running
	Paused
	Saved
	Stopped
	Stopping
	Starting
	Reloading
	Error
	Timeout
)

var states = []string{
	"None",
	"Running",
	"Paused",
	"Saved",
	"Stopped",
	"Stopping",
	"Starting",
	"Reloading",
	"Error",
	"Timeout",
}

//To string of State
func (s State) String() string {
	if int(s) >= 0 && int(s) < len(states) {
		return states[int(s)]
	}
	return "None"
}
