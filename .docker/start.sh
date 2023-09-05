#!/bin/sh

# We'll use this script to manage starting and stopping this container gracefully.
# It only takes up about 00.01 CPU % allotted to the container, you can verify
# by running `docker stats` after you start a container that uses this as
# as the CMD.

set -e

shutd () {
    printf "%s" "Shutting down the container gracefully..."

    # You can run clean commands here!

    echo "done"
}

trap 'shutd' TERM

echo "Starting up"

# Run non-blocking commands here
go mod tidy
go mod vendor

echo "Ready."

# This will keep the container running, it has very low CPU usage.
# When this process receives a signal to be stopped, shutd function will run.
while :; do :; done & kill -STOP $! && wait $!
