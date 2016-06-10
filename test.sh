#!/bin/bash
set -e; set -o pipefail

function fatal() {
	echo $@>&2
	exit 1
}

[ -n "$GOPATH" ] || fatal "ERROR: Missing GOPATH env var."
[ -e /tmp/config.json ] || fatal "ERROR: Put test config.json in /tmp."

./build.sh
$GOPATH/bin/provisioner-server /tmp/provisioner.sock /tmp/config.json &
server_pid=$!

# Give some time for startup.
sleep 1

$GOPATH/bin/provisioner-test-client /tmp/provisioner.sock
kill $server_pid
