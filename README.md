# Simple worker for ArchCI

## Introduction

Simple-worker is the easy-to-deploy worker to run continues integreation tests.

It pulls tasks from ArchCI service and run the test within docker containers. We can use any management tool for docker to control these task and it's much more efficient than Jenkins and TravisCI.

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

Or you can run simple-worker with `docker run -d --net=host archci/simple-work`.