FROM golang:1.4

ADD . /go/src/github.com/jmcarbo/docker-syslog
RUN go get github.com/jmcarbo/docker-syslog/...
RUN go install github.com/jmcarbo/docker-syslog

EXPOSE 1514/udp

ENTRYPOINT [ "/go/bin/docker-syslog" ]
