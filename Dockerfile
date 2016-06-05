FROM golang
COPY . /go/src/github.com/kamilmac/cubesdb
WORKDIR /go/src/github.com/kamilmac/cubesdb
RUN go get golang.org/x/net/context
RUN go get github.com/satori/go.uuid
RUN go get github.com/boltdb/bolt
RUN go get goji.io
RUN go build main.go
CMD ["./main"]