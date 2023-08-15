# LogSync

Fetch and filter events from Swiss Learning Hub. Report events to syslog server in CEF format.

## Usage

See examples below on how to run and test logsync:

```shell
# Check version
logsync -v
#> logsync version 1.0.0

# Show help for app and command
logsync -h
logsync run -h

# Run in dry-run mode first to check functionality/plausibility
# No events will be forwareded if -d|--dry-run is set
logsync run -d
logsync run -d -c /path/to/my.json

# If satisfied...
logsync run
logsync run -c /path/to/my.json
```

## Configuration

If not given else via CLI flags, logsync will try to find and read a `logsync.json` file in

1. Current working directory
2. Directory of binary

See `logsync.json` schema [here](logsync.dist.json). For more details also see [detections & reporters](detections.md).

## Local Logging (Optional)

By default, output is written to StdOut. Configuration setting `logfile` can be configured to also write to file. File
itself is getting created if not already exists (permissions `0644`). The directory must exist and will not be created
automatically.

## Filter

Filter are used to limit queried events from SLH. The `days` parameter is mandatory:

| Attribute | Type         | Info                                                    |
|-----------|:-------------|---------------------------------------------------------|
| `days`    | `<int>`      | Days to fetch events from (`1-90`)                      |
| `type`    | `[]<string>` | Optional: Limit events to this array of types.          |
| `max`     | `<int>`      | Optional: Maximum entries to retrieve (default: 999999) |

## Logging Facility

Define in syslog configuration `facility` (suggestion: `32`):

```
LOG_KERN     = 0
LOG_USER     = 8
LOG_MAIL     = 16
LOG_DAEMON   = 24
LOG_AUTH     = 32
LOG_SYSLOG   = 40
LOG_LPR      = 48
LOG_NEWS     = 56
LOG_UUCP     = 64
LOG_CRON     = 72
LOG_AUTHPRIV = 80
LOG_FTP      = 88
```
