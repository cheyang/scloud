// interfaces.go
package msg

type Sender interface {
	Send(entry interface{})
}

type Receiver interface {
	Recieve() interface{}
}
