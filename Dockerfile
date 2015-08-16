from busybox:latest
maintainer Cory Buecker <email@corybuecker.com>

env GOROOT /go
add https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /go/lib/time/zoneinfo.zip
add https://gist.githubusercontent.com/corybuecker/50f922f6ffc52d291f00/raw/dc2854fbb307dab11191b5d63288f7ffcf59e885/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

add trade-fetcher /usr/bin/trade-fetcher
entrypoint ["/usr/bin/trade-fetcher"]
