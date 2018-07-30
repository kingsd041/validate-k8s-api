FROM golang:1.10.3

RUN go get github.com/rancher/longhorn-manager && \
    go get github.com/coreos/etcd

COPY ./mytest.go /go/src/github.com/rancher/longhorn-manager/mytest.go

COPY ./test /go/src/github.com/rancher/longhorn-manager/test

WORKDIR /go/src/github.com/rancher/longhorn-manager

ENTRYPOINT ["./test"]