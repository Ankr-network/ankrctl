FROM golang:1.10-alpine3.8 as builder
ARG URL_BRANCH
ARG TENDERMINT_URL
RUN apk update && \
    apk add git && \
    apk add --update bash && \
    apk add openssh
RUN go get github.com/golang/dep/cmd/dep

COPY id_rsa /root/.ssh/
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN chmod go-w /root
RUN chmod 700 /root/.ssh
RUN chmod 600 /root/.ssh/id_rsa

WORKDIR $GOPATH/src/github.com/Ankr-network/dccn-cli/
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only
COPY . $GOPATH/src/github.com/Ankr-network/dccn-cli/

RUN echo ${URL_BRANCH}
RUN echo ${TENDERMINT_URL}
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -a \
    -installsuffix cgo \
    -ldflags="-w -s -X github.com/Ankr-network/dccn-cli/commands.clientURL=${URL_BRANCH} -X github.com/Ankr-network/dccn-cli/commands.tendermintURL=${TENDERMINT_URL}" \
    -o /go/bin/ankrctl \
    cmd/ankrctl/main.go

FROM alpine:3.7
COPY --from=builder /go/bin/ankrctl /bin/ankrctl

