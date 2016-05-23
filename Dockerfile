FROM golang
COPY . /go/src/github.com/kamilmac/cubesdb
WORKDIR /go/src/github.com/kamilmac/cubesdb
RUN go build main.go
CMD ["./main"]