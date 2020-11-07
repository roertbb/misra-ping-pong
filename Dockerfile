FROM golang:1.14.2-alpine

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

RUN apk update && apk add zeromq-dev alpine-sdk
WORKDIR /home

COPY . .
RUN make deps
RUN make build

ENV PING_LOSS_PROB=0 \
    PONG_LOSS_PROB=0 \
    ADDRESSES="127.0.0.1:3001 127.0.0.1:3002 127.0.0.1:3003"

CMD ./bin/main \
    -ping-loss-prob=$PING_LOSS_PROB \
    -pong-loss-prob=$PONG_LOSS_PROB \
    $ADDRESSES