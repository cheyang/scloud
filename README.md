# scloud
Setup in cloud with golang


### Update GoDeps

```
cd $GOPATH/src/github/cheyang/scloud/cmd/scloud
godep restore
rm -rf Godeps
go get github.com/cheyang/scloud/cmd/scloud
cd $GOPATH/src/github.com/cheyang/scloud
godep save ./...
```


### Develop unit test

```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
cd path/to/package/you/want/to/test

ginkgo bootstrap # set up a new ginkgo suite
ginkgo generate  # will create a sample test file.  edit this file and add your tests then...

go test # to run your tests

ginkgo  # also runs your tests
```