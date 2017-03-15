# Icingabeat CHANGELOG

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
