#!/bin/bash
set -e

/usr/bin/dockerlaunch /usr/bin/docker -d -s overlay &

/go/src/github.com/ArchCI/simple-worker/simple-worker
