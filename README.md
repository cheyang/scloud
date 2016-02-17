# scloud
Setup in cloud with golang


### Update GoDeps

```
cd $GOPATH/src/github/cheyang/scloud/cmd/scloud
godep restore
rm -rf Godeps
go get github.com/cheyang/scloud/cmd/scloud
godep save ./...
```
