FROM golang:alpine as builder

RUN apk update && apk add git && apk add --update bash && apk add openssh

COPY . $GOPATH/src/github.com/Ankr-network/dccn-cli/
COPY id_rsa /root/.ssh/
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN chmod go-w /root
RUN chmod 700 /root/.ssh
RUN chmod 600 /root/.ssh/id_rsa
WORKDIR $GOPATH/src/github.com/Ankr-network/dccn-cli/
RUN git clone -b feat/swdev-79-dccncli git@github.com:Ankr-network/godo.git $GOPATH/src/github.com/Ankr-network/godo --config core.autocrlf=input
RUN git clone -b feature/78-ankr-hub git@github.com:Ankr-network/dccn-hub.git $GOPATH/src/github.com/Ankr-network/dccn-hub --config core.autocrlf=input

RUN go get $GOPATH/src/github.com/Ankr-network/dccn-cli/cmd/dccncli/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/ankr $GOPATH/src/github.com/Ankr-network/dccn-cli/cmd/dccncli/main.go

FROM scratch
COPY --from=builder /go/bin/ankr /go/bin/ankr
ENTRYPOINT ["/go/bin/ankr"]