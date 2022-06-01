# Icingabeat CHANGELOG

## v7.17.4

### Features
* Update libbeat to version 7.14.2

### Breaking Changes
* Dashboards now must be importaed manually using Kibana

## v7.14.2

### Features
* Update libbeat to version 7.14.2

## v7.5.2

### Features
* Update libbeat to version 7.5.2

## v7.4.2

### Features
* Update libbeat to version 7.4.2

## v6.5.4

### Features
* Update libbeat to version 6.5.4
* Move all field names to 'icinga' namespace

### Bugs
* Prevent usage of reserved keywords

## v6.3.3

### Features
* Update libbeat to version 6.3.3

### Bugs
* Remove `zones` key from statuspoller. This key may become to big to process.
* Catch 404 return codes
* Update dashboard directory schema so `icingabeat setup` works out of the box

## v6.1.1

### Features
* Update libbeat to version 6.1.1
* Add setting to add custom CAs for SSL verification

### Bugs
* Close connections properly on failed authentication

## v5.6.6

### Features
* Update libbeat to version 5.6.6

## v1.1.1

### Bugs
* Add missing fields for statuspoller

### Features
* Update libbeat to version 5.3.2

## v1.1.0

### Features
* Update libbeat to version 5.3.0

## v1.0.0

### Features
* Add dashboards
* Cleanup logging

## v0.2.0

### Features
* Add statuspoller
* Add option to skip ssl verificatoin
* Move retry_interval setting to eventstream scope
* Add field types and documentation
* Use libbeat 5.1.2

### Bugs
* Fix packaging scripts
* Process events
  * Capture exceptions
  * Rename fields
  * Create proper Timestamps

## v0.1.0

### Features
* Configure connection settings
* Configure eventstream types
* Configure eventstream filters
* Subscribe to eventstream
* Decument fields for eventstream types:
  * CheckResult
  * StateChange
  * Notification
  * AcknowledgementSet
  * AcknowledgementCleared
