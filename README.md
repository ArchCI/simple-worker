# Simple worker for ArchCI

## Introduction

Simple-worker is the easy-to-deploy worker to run continues integreation tasks for ArchCI.

It pulls tasks from ArchCI service and run the test within docker containers. We can use any docker management tool for our tasks and it's much more efficient than Jenkins or TravisCI.

## Install

```
go get github.com/ArchCI/simple-worker
```

Or build from source.

```
cd ArchCI/simple-worker/
go build
```

## Usage

```
./simple-worker
```

Simply running the binary will start the agent to get task to test. You can setup the configuration with `worker.yml`.

## Docker container

ArchCI relies on docker to run the tests. Make sure that docker daemon is running on your server.

Or you can run simple-worker within container by `docker run -d --net=host archci/simple-work`.