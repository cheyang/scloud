package msg_test

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	//	"github.com/cheyang/scloud/drivers/softlayer"
	"github.com/cheyang/scloud/pkg/drivers"
	"github.com/cheyang/scloud/pkg/host"
	. "github.com/cheyang/scloud/pkg/msg"
	"github.com/davecgh/go-spew/spew"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Msg", func() {
	var (
		hosts []*host.Host
		num   int
	)

	Context("#Pub for specified number", func() {

		BeforeEach(func() {

			num = 2

			hosts = make([]*host.Host, num)

			for i := 0; i < num; i++ {
				hosts[i] = &host.Host{Name: strconv.Itoa(i),
					Driver: &drivers.BaseDriver{IPAddress: strconv.Itoa(i),
						MachineName: strconv.Itoa(i)},
				}
				fmt.Fprintf(os.Stdout, "exec method GetMachineName for %s\n", hosts[i].Driver.GetMachineName())
			}
		})

		runtime.GOMAXPROCS(runtime.NumCPU())
		It("Test members", func() {
			for i := 0; i < num; i++ {
				fmt.Fprintf(os.Stdout, "receive  %d %s", i, hosts[i])
			}

			Expect(len(hosts)).ToNot(BeNil())

			queue := NewQueue(num)

			defer queue.Close()

			for i := 0; i < num; i++ {
				go func(n int) {
					fmt.Fprintf(os.Stdout, "call %d\n", n)
					queue.Send(hosts[n])
				}(i)
			}

			for i := 0; i < num; i++ {
				entry := queue.Recieve()
				//				fmt.Fprintf(os.Stdout, "receive %s\n", entry)
				spew.Printf("entry =%#+v\n", entry)

				if h, ok := entry.(*host.Host); ok {
					fmt.Fprintf(os.Stdout, "exec method GetMachineName for %s\n", h.Driver.GetMachineName())
				}
			}
		})

	})
})
