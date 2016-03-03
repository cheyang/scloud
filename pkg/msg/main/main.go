// main.go
package main

import (
	"fmt"
	"os"
	//	"runtime"
	"strconv"
	//	"github.com/cheyang/scloud/drivers/softlayer"
	"github.com/cheyang/scloud/pkg/drivers"
	"github.com/cheyang/scloud/pkg/host"
	. "github.com/cheyang/scloud/pkg/msg"
)

func main() {
	var (
		hosts []*host.Host
		num   int
	)

	num = 5

	hosts = make([]*host.Host, num)

	for i := 0; i < num; i++ {
		hosts[i] = &host.Host{Name: strconv.Itoa(i),
			Driver: &drivers.BaseDriver{IPAddress: strconv.Itoa(i)}}

	}

	for i := 0; i < num; i++ {
		fmt.Fprintf(os.Stdout, "Print  %d %s\n", i, hosts[i])
	}

	//	Expect(len(hosts)).ToNot(BeNil())

	queue := NewQueue(num)

	defer queue.Close()

	for i := 0; i < num; i++ {
		go func(n int) {
			fmt.Fprintf(os.Stdout, "call %d", n)
			queue.Send(hosts[n])
		}(i)
	}

	for i := 0; i < num; i++ {
		entry := queue.Recieve()
		fmt.Fprintf(os.Stdout, "receive %s", entry)
	}
}
