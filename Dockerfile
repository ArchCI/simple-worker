FROM golang:1.4
MAINTAINER tobe tobeg3oolge@gmail.com

# Support docker in docker
ADD rancher_docker.tar /

# Install tools
RUN apt-get update -y
RUN apt-get remove -y git
RUN apt-get install -y git
RUN go get github.com/tools/godep

# Build simple-worker
ADD . /go/src/github.com/ArchCI/simple-worker
WORKDIR /go/src/github.com/ArchCI/simple-worker
RUN godep go build -ldflags "-X main.GitVersion `git rev-parse HEAD` -X main.BuildTime `date -u '+%Y-%m-%d_%I:%M:%S'`"

VOLUME /var/lib/docker

CMD ./entrypoint.sh