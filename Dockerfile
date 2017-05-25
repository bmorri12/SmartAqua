FROM golang:1.5.2

MAINTAINER xfocus xfocus3@gmail.com

ADD . /go/src/github.com/bmorri12/SmartAqu

RUN go get github.com/bmorri12/SmartAqu/...