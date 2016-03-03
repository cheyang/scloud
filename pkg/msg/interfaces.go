// interfaces.go
package msg

type Sender interface {
	Send(entry interface{})
	SetDone()
}

type Receiver interface {
	Recieve() interface{}
	IsDone() bool // Check if the sender finished sending msg

}
