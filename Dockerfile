FROM golang

ADD . /go/github.com/davidahouse/socket-redis

RUN go get github.com/go-redis/redis && go get github.com/gorilla/websocket && go install socket-redis

ENTRYPOINT /go/bin/socket-redis

EXPOSE 8081
