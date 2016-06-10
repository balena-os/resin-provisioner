#!/bin/bash
set -e; set -o pipefail

go get github.com/gorilla/mux
go get github.com/coreos/go-systemd/dbus
go install ./...
