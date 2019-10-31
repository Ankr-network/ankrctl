FROM golang:1.13-alpine as builder
ARG URL_BRANCH
ARG TENDERMINT_URL
ARG TENDERMINT_PORT
ARG GITHUB_USER
ARG GITHUB_TOKEN
ARG GOPROXY
ENV GOPROXY=${GOPROXY}
ENV GOPRIVATE=github.com/Ankr-network
RUN echo "machine github.com login ${GITHUB_USER} password ${GITHUB_TOKEN}" > ~/.netrc

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
    main.go

FROM alpine:3.7
RUN  apk update && \
     apk add libc6-compat && \
     apk add ca-certificates
COPY --from=builder /go/bin/ankrctl /bin/ankrctl

