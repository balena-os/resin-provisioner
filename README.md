# resin-provisioner

This is a resin supervisor component designed to allow provisioning of a resin
OS device against resin servers.

## Usage

The core functionality is provided by the `provisioner` package, and 2 utility
commands are provided:

### provisioner-server

Offers provisioner functionality as an HTTP API via the specified socket file.

```
$ provisioner-server [socket path] [config.json path]
```

### provisioner-simple-client

This is a simple provisioning tool. To query the provisioned state use:

```
$ provisioner-simple-client [config path]
```

To execute a provisioning use:

```
$provisioner-simple-client [config path] [user id] [app id] [api key]
```

It uses the following environment variables:

* `RESIN_PUBNUB_SUBSCRIBE_KEY` - Sets pubnub subscribe key if not specified
  already in config.json.
* `RESIN_PUBNUB_PUBLISH_KEY` - Sets pubnub publish key if not specified already
  in config.json.
* `RESIN_MIXPANEL_TOKEN` - Sets mixpanel token if not specified already in
  config.json.
* `RESIN_DOMAIN_OVERRIDE` - Overrides the default domain (`resin.io`) and sets
  it to this value.
