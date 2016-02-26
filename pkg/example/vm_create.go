// vm_create.go
package main

import (
	"fmt"

	lib "github.com/cheyang/scloud/pkg"
	"github.com/cheyang/scloud/pkg/host"
	"github.com/cheyang/scloud/pkg/host/errs"
)

func main() {
	fmt.Println("Hello World!")

	store := lib.GetDefaultStore("test_cluster")

	hostname := "myhost"

	validateName := host.validHostName(hostname)

	if !validateName {
		fmt.Printf("Error creating machine: %s", errs.ErrInvalidHostname)
		return
	}

}
