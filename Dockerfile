FROM golang:1.4
MAINTAINER tobe tobeg3oolge@gmail.com

# Manage dependency
RUN go get github.com/tools/godep

# Support docker in docker
ADD rancher_docker.tar /

# Build simple-worker
ADD . /go/src/github.com/ArchCI/simple-worker
WORKDIR /go/src/github.com/ArchCI/simple-worker
RUN godep go build

VOLUME /var/lib/docker

CMD ./entrypoint.sh