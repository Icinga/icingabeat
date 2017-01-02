[![Build Status](https://travis-ci.org/Icinga/icingabeat.svg?branch=master)](https://travis-ci.org/Icinga/icingabeat)

# Icingabeat

Icingabeat is a Beat to fetch data from the Icinga 2 API. Icingabeat has two
modes:

**Eventstream:** Receive an eventstream from the Icinga 2 API. Events are
checkresults, notifications and many other types. See below for details. There
is no polling involved when receiving an eventstream.

Example usage:
* Correlate monitoring data with logging information
* Create dashboards in Kibana with metrics collected by Icinga 2
* Monitor notifications sent by Icinga 2

**StatusPoller:** The Icinga 2 API exports lots of information about the state
of the Icinga daemon. These information can be polled periodically.

Example usage:
* Metrics of Icinga 2 performance
* Metrics about each enabled Icinga 2 feature
* Information about zones and endpoints


## Installation
There are no packages available yet. Please build and run Icingabeat manually.

## Configuration
Before starting Icingabeat make sure you have set up your configuration properly.
There are several general settings and some that are specific to the modes
included.

### Connection
These settings apply to both modes. They define the API endpoint to which
Icingabeat connects to.

#### `host`
Hostname of your Icinga 2 API. This can be either an IP address or domain.
Defaults to `localhost`

#### `port`
Defaults to `5665`

#### `user`
Username for the API connection. In your Icinga 2 configuration, this is a user
of the object type `apiuser`. Make sure that it has sufficient permissions to
read the data you want to collect. Learn more about `apiuser` permissions in the
[Icinga 2 docs](https://docs.icinga.com/icinga2/latest/doc/module/icinga2/chapter/icinga2-api#icinga2-api-permissions).

This user is allowed to receive all types of events and to query the status API:
```c++
object ApiUser "icinga" {
  password: "icinga"
  permissions = ["events/*", "status/query"]
}
```

This user is only allowed to receive events of the type `CheckResult`
```c++
object ApiUser "icinga" {
  password: "icinga"
  permissions = ["events/CheckResult"]
}
```

#### `password`
Defaults to `icinga`

#### `retry_interval`
Instead of stopping on connection loss, Icingabeat will try to reconnect to the
API periodically. Defaults to `10s`

### Eventstream
These settings are eventsream specific, they apply only to the eventstream. To
disable evenstream completely, comment out the whole section.

#### `types`
Decide which events you want to receive from the event stream. The following
event stream types are available, at least one must be set:

* CheckResult
* StateChange
* Notification
* AcknowledgementSet
* AcknowledgementCleared
* CommentAdded
* CommentRemoved
* DowntimeAdded
* DowntimeRemoved
* DowntimeStarted
* DowntimeTriggered

#### `filter`
Event streams can be filtered by attributes using the prefix `event.` By default
no filter is set.


Only checkresults with exit status 2:
```yaml
filter: "event.check_result.exit_status==2"
```

Example for the CheckResult type with the service matching the string pattern
`mysql*`:
```yaml
filter: 'match("mysql*", event.service)'
```

### StatusPoller
StatusPoller is not implemented yet.

## Run
To run Icingabeat with debugging output enabled, run:

```bash
./icingabeat -c icingabeat.yml -e -d "*"
```

Icingabeat exports a bunch of fields. Have a look to the
[fields.asciidoc](docs/fields.asciidoc) for details.

## Building and running manually

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Clone

To clone Icingabeat from the git repository, run the following commands:

```shell
mkdir -p ${GOPATH}/github.com/icinga
cd ${GOPATH}/github.com/icinga
git clone https://github.com/icinga/icingabeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build
Ensure that this folder is at the following location:
`${GOPATH}/github.com/icinga`

To build the binary for Icingabeat run the command below. This will generate a
binary in the same directory with the name icingabeat.

```shell
make
```

### Run
To run Icingabeat with debugging output enabled, run:

```shell
./icingabeat -c icingabeat.yml -e -d "*"
```

### Test

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

### Update

Each beat has a template for the mapping in elasticsearch and a documentation
for the fields which is automatically generated based on `etc/fields.yml`.
To generate etc/icingabeat.template.json and etc/icingabeat.asciidoc

```shell
make update
```

### Cleanup

To clean  Icingabeat source code, run the following commands:

```shell
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```shell
make clean
```

## Packaging

The beat frameworks provides tools to crosscompile and package your beat for
different platforms. This requires [docker](https://www.docker.com/) and
vendoring as described above. To build packages of your beat, run the following
command:

```shell
make package
```

This will fetch and create all images required for the build process. The hole
process to finish can take several minutes.
