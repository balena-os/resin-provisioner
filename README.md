# resin-provisioner

This tool is a resin component for converting an unmanaged resinOS device into a
managed resinOS device. This is achieved by provisioning the device against the
resin.io servers.

The core functionality is provided by the `resin-provision` binary which
includes two modes of interaction: interactive mode and oneshot mode.

## Interactive mode
This mode is interactive and allows the user to create new accounts and
applications. It is used by the `local promote` command in the [Resin
CLI](https://github.com/resin-io/resin-cli)

## Oneshot mode
This mode allows the user to provision a resinOS device with a single command.
It is designed to be used within external applications or as part of an
automated process.
