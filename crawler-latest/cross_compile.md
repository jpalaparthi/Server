# Windows 64 bit architecture compilation

env GOOS="windows" GOARCH="amd64" CGO_ENABLED="1" GOPATH="/Users/mac/workspace/go-apps" go build "src/crawler/crawler.go"

env GOOS="windows" GOARCH="386" CGO_ENABLED="1" GOPATH="/Users/mac/workspace/go-apps" go build "src/crawler/crawler.go"