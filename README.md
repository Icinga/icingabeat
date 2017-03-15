[![Build Status](https://travis-ci.org/Icinga/icingabeat.svg?branch=master)](https://travis-ci.org/Icinga/icingabeat)

# Icingabeat

> The Beats are lightweight data shippers, written in Go, that you install on
> your servers to capture all sorts of operational data (think of logs,
> metrics, or network packet data). The Beats send the operational data to
> Elasticsearch, either directly or via Logstash, so it can be visualized with
> Kibana.

Icingabeat is an [Elastic Beat](https://www.elastic.co/products/beats) that
fetches data from the Icinga 2 API and sends it either directly to Elasticsearch
or Logstash. This Beat supports two modes:

## Eventstream

Receive an eventstream from the Icinga 2 API. This stream includes events such
as checkresults, notifications, downtimes, acknowledgemts and many other types.
See below for details. There is no polling involved when receiving an
eventstream.

Example usage:
* Correlate monitoring data with logging information
* Monitor notifications sent by Icinga 2

## Statuspoller

The Icinga 2 API exports a lot of information about the state of the Icinga
daemon. Icingabeat can poll these information periodically.

Example usage:
* Visualize metrics of the Icinga 2 daemon
* Get insights how each enable Icinga 2 feature performs
* Information about zones and endpoints


### Installation
Download and install your package from the
[latest release](https://github.com/Icinga/icingabeat/releases/latest) page.

### Configuration
Configuration of Icingabeat is splitted into 3 sections: General, Evenstream and
Statuspoller. On Linux configuration files are located at `/etc/icingabeat`

#### General
Settings in this section apply to both modes.

##### `host`
Hostname of Icinga 2 API. This can be either an IP address or domain.
Defaults to `localhost`

##### `port`
Defaults to `5665`

##### `user`
Username to be used for the API connection. You need to create this user in your Icinga 2 configuration. Make sure that it has sufficient permissions to read the
data you want to collect.

Here is an example how an API user in your Icinga 2 configuration:

```c++
object ApiUser "icinga" {
  password: "icinga"
  permissions = ["events/*", "status/query"]
}
```


Learn more about the `ApiUser` and its permissions in the
[Icinga 2 docs](https://docs.icinga.com/icinga2/latest/doc/module/icinga2/chapter/icinga2-api#icinga2-api-permissions).

##### `password`
Defaults to `icinga`

##### `skip_ssl_verify`
Skip verification of SSL certificates. Defaults to `false`

#### Eventstream
Settings in this section apply to the eventstream mode. To disable the
eventstream completely, comment out the section.

##### `types`
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
  types:
    - CheckResult
    - StateChange
    - Notification
    - AcknowledgementSet
    - AcknowledgementCleared
```

##### `filter`
Additionally to selecting the types of events, you can filter them  by
attributes using the prefix `event.` and  By default no filter is set.

###### Examples

Only checkresults with the exit code 2:
```yaml
  filter: "event.check_result.exit_status==2"
```

Only checkreults of services that match `mysql*`:
```yaml
  filter: 'match("mysql*", event.service)'
```

##### `retry_interval`
On a connection loss Icingabeat will try to reconnect to the API periodically.
This setting defines the interval for connection retries. Defaults to `10s`

#### Statuspoller
Settings of this section apply to the statuspoller mode.

##### `interval`
Interval at which the status API is called. Set to `0` to disable polling.
Defaults to `60s`

### Run

To start Icingabeat, use your operating systems default commands. On Linux this
should be one of these commands, depending on the distribution you are using:

* `service icingabeat start`
* `systemctl icingabeat start` or
* `/etc/init.d/icingabeat start`

## Fields
Icingabeat exports a bunch of fields. Have a look to the
[fields.asciidoc](docs/fields.asciidoc) for details.

## Development

### Building and running manually

#### Requirements

* [Golang](https://golang.org/dl/) 1.7

#### Clone

To clone Icingabeat from the git repository, run the following commands:

```shell
mkdir -p ${GOPATH}/github.com/icinga
cd ${GOPATH}/github.com/icinga
git clone https://github.com/icinga/icingabeat
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

#### Build
Ensure that this folder is at the following location:
`${GOPATH}/github.com/icinga`

To build the binary for Icingabeat run the command below. This will generate a
binary in the same directory with the name icingabeat.

```shell
make
```

#### Run
To run Icingabeat with debugging output enabled, run:

```shell
./icingabeat -c icingabeat.yml -e -d "*"
```

#### Test

To test Icingabeat, run the following command:

```shell
make testsuite
```

alternatively:
```shell
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

#### Update

Each beat has a template for the mapping in elasticsearch and a documentation
for the fields which is automatically generated based on `etc/fields.yml`.
To generate etc/icingabeat.template.json and etc/icingabeat.asciidoc

```shell
make update
```

#### Cleanup

To clean  Icingabeat source code, run the following commands:

```shell
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```shell
make clean
```

### Packaging

The beat frameworks provides tools to crosscompile and package your beat for
different platforms. This requires [docker](https://www.docker.com/) and
vendoring as described above. To build packages of your beat, run the following
command:

```shell
make package
```

This will fetch and create all images required for the build process. The hole
process to finish can take several minutes.

To disable snapshot packages or build specific packages, set the following
environment variables:

```shell
export SNAPSHOT=false
export TARGETS="\"linux/amd64 linux/386\""
export PACKAGES=icingabeat/deb
make package
```
