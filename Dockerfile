FROM golang:1.4
MAINTAINER tobe tobeg3oolge@gmail.com

RUN apt-get update -y

RUN apt-get install -y git
RUN apt-get install -y docker.io

ADD . /go/simple-worker
WORKDIR /go/simple-worker

RUN go get
RUN go build

CMD /bin/bash