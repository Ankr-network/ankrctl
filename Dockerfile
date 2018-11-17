FROM golang:alpine as builder

RUN apk update && apk add git && apk add --update bash

COPY . $GOPATH/src/github.com/Ankr-network/dccn-cli/
COPY id_rsa /root/.ssh/
WORKDIR $GOPATH/src/github.com/Ankr-network/dccn-cli/
RUN git clone -b feat/swdev-79-dccncli https://github.com/Ankr-network/godo.git $GOPATH/src/github.com/Ankr-network/godo --config core.autocrlf=input
RUN git clone -b feature/78-ankr-hub https://github.com/Ankr-network/dccn-hub.git $GOPATH/src/github.com/Ankr-network/dccn-hub --config core.autocrlf=input

RUN go get $GOPATH/src/github.com/Ankr-network/dccn-cli/cmd/ankr/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/ankr $GOPATH/src/github.com/Ankr-network/dccn-cli/cmd/ankr/main.go

FROM scratch
COPY --from=builder /go/bin/ankr /go/bin/ankr
ENTRYPOINT ["/go/bin/ankr"]