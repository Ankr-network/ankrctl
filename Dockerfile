FROM golang:1.12-alpine as builder
ARG URL_BRANCH
ARG TENDERMINT_URL
ARG TENDERMINT_PORT
RUN apk update && \
    apk add git && \
    apk add --update bash && \
    apk add openssh
#RUN go get github.com/golang/dep/cmd/dep

COPY id_rsa /root/.ssh/
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN chmod go-w /root
RUN chmod 700 /root/.ssh
RUN chmod 600 /root/.ssh/id_rsa

WORKDIR $GOPATH/src/github.com/Ankr-network/ankrctl/
#COPY Gopkg.toml Gopkg.lock ./
#RUN dep ensure -vendor-only
COPY . $GOPATH/src/github.com/Ankr-network/ankrctl/

RUN export GO111MODULE=on
RUN go mod download

RUN echo ${URL_BRANCH}
RUN echo ${TENDERMINT_URL}
RUN echo ${TENDERMINT_PORT}
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -a \
    -installsuffix cgo \
    -ldflags="-w -s -X github.com/Ankr-network/ankrctl/commands.clientURL=${URL_BRANCH} -X github.com/Ankr-network/ankrctl/commands.tendermintURL=${TENDERMINT_URL} -X github.com/Ankr-network/ankrctl/commands.tendermintPort=${TENDERMINT_PORT}" \
    -o /go/bin/ankrctl \
    cmd/ankrctl

FROM alpine:3.7
RUN  apk update && \
     apk add libc6-compat && \
     apk add ca-certificates
COPY --from=builder /go/bin/ankrctl /bin/ankrctl

