# Configuration
Configuration of Icingabeat is split into 3 sections: Connection, Evenstream and
Statuspoller. On Linux configuration files are located at `/etc/icingabeat`

## Connection
Settings in this section apply to both modes.

### `host`
Hostname of Icinga 2 API. This can be either an IP address or domain.
Defaults to `localhost`

### `port`
Defaults to `5665`

### `user`
Username to be used for the API connection. You need to create this user in your Icinga 2 configuration. Make sure that it has sufficient permissions to read the
data you want to collect.

Here is an example of an API user in your Icinga 2 configuration:

``` c++
object ApiUser "icinga" {
  password = "icinga"
  permissions = ["events/*", "status/query"]
}
```

Learn more about the `ApiUser` and its permissions in the
[Icinga 2 docs](https://docs.icinga.com/icinga2/latest/doc/module/icinga2/chapter/icinga2-api#icinga2-api-permissions).

### `password`
Defaults to `icinga`


### `ssl.verify`
Configure SSL verification. If `false` is configured, all server hosts and
certificates will be accepted. In this mode, SSL based connections are
susceptible to man-in-the-middle attacks. Use only for testing. Default is
`true`.

### `ssl.certificate_authorities`
List of root certificates for HTTPS server verifications

Example:
```
ssl.certificate_authorities: ["/etc/pki/root/ca.pem"]
```

## Eventstream
Settings in this section apply to the eventstream mode. To disable the
eventstream completely, comment out the section.

### `eventstream.types`
You can select which particular Icinga 2 events you want to receive and store.
The following types are available, you must set at least one:

* `CheckResult`
* `StateChange`
* `Notification`
* `AcknowledgementSet`
* `AcknowledgementCleared`
* `CommentAdded`
* `CommentRemoved`
* `DowntimeAdded`
* `DowntimeRemoved`
* `DowntimeStarted`
* `DowntimeTriggered`

To set multiple types, do the following:

```yaml
  eventstream.types:
    - CheckResult
    - StateChange
    - Notification
    - AcknowledgementSet
    - AcknowledgementCleared
```

### `eventstream.filter`
In addition to selecting the types of events, you can filter them by
attributes using the prefix `event.`. By default no filter is set.

###### Examples

Only check results with the exit code 2:
```yaml
  eventstream.filter: "event.check_result.exit_status==2"
```

Only check results of services that match `mysql*`:
```yaml
  eventstream.filter: 'match("mysql*", event.service)'
```

### `eventstream.retry_interval`
On a connection loss Icingabeat will try to reconnect to the API periodically.
This setting defines the interval for connection retries. Defaults to `10s`

## Statuspoller
Settings of this section apply to the statuspoller mode.

### `statuspoller.interval`
Interval at which the status API is called. Set to `0` to disable polling.
Defaults to `60s`

## Fields
Icingabeat exports a bunch of fields. Have a look to the
[fields.asciidoc](https://github.com/Icinga/icingabeat/blob/master/docs/fields.asciidoc) for details.
