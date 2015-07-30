#!/bin/bash
set -e

# Use overlay by defulat, if set AUFS=true just use aufs
if [ -z "$AUFS" ]; then
    echo "Use defualt overlay filesystem"
    /usr/bin/dockerlaunch /usr/bin/docker -d -s overlay &
else
    echo "Use aufs filesystem"
    /usr/bin/dockerlaunch /usr/bin/docker -d -s aufs &
fi

sleep 5
/go/src/github.com/ArchCI/simple-worker/simple-worker
