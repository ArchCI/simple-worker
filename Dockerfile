FROM golang:1.4
MAINTAINER tobe tobeg3oolge@gmail.com

RUN apt-get update -y

RUN apt-get install -y git

ADD rancher_docker.tar /

ADD . /go/src/github.com/ArchCI/simple-worker
WORKDIR /go/src/github.com/ArchCI/simple-worker

RUN go get
RUN go build

VOLUME /var/lib/docker

CMD /bin/bash