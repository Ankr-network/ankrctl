FROM golang:1.10-alpine3.8 as builder
ARG URL_BRANCH
RUN apk update && \
    apk add git && \
    apk add --update bash && \
    apk add openssh

WORKDIR /ankrctl
COPY . /ankrctl

RUN GOPROXY=https://goproxy.cn CGO_ENABLED=0 \
    go build -a \
    -ldflags="-w -s -X github.com/Ankr-network/ankrctl/commands.clientURL=${URL_BRANCH} -X github.com/Ankr-network/ankrctl/commands.tendermintURL=${TENDERMINT_URL} -X github.com/Ankr-network/ankrctl/commands.tendermintPort=${TENDERMINT_PORT}" \
    -o /go/bin/ankrctl \
    cmd/ankrctl/main.go

FROM scratch
COPY --from=builder /go/bin/ankrctl /go/bin/ankrctl
ENTRYPOINT ["/go/bin/ankrctl"]
