FROM golang:1.12-alpine as builder
ARG URL_BRANCH
ARG TENDERMINT_URL
ARG TENDERMINT_PORT
RUN apk update && \
    apk add git && \
    apk add --update bash && \
    apk add openssh
#RUN go get github.com/golang/dep/cmd/dep

WORKDIR /ankrctl
COPY . /ankrctl

RUN GOPROXY=https://goproxy.cn CGO_ENABLED=0 \
    go build -a \
    -ldflags="-w -s -X github.com/Ankr-network/ankrctl/commands.clientURL=${URL_BRANCH} -X github.com/Ankr-network/ankrctl/commands.tendermintURL=${TENDERMINT_URL} -X github.com/Ankr-network/ankrctl/commands.tendermintPort=${TENDERMINT_PORT}" \
    -o /go/bin/ankrctl \
    cmd/ankrctl/main.go

FROM alpine:3.7
RUN  apk update && \
     apk add libc6-compat && \
     apk add ca-certificates
COPY --from=builder /go/bin/ankrctl /bin/ankrctl

