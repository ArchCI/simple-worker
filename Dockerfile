FROM ubuntu:15.04

RUN apt-get update -y

RUN apt-get install -y git
RUN apt-get install -y docker.io

ADD . /simple-worker
WORKDIR /simple-worker

CMD /bin/bash