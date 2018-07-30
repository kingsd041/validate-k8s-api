FROM golang:1.10.3

RUN mkdir -p /go/src/github.com/rancher/longhorn-manager/

WORKDIR /go/src/github.com/rancher/longhorn-manager

RUN git clone https://github.com/rancher/longhorn-manager.git

COPY . /go/src/github.com/rancher/longhorn-manager/

CMD ["./test"]