# Detections & Reporters

Using multiple detections and reporters allows to filter different kind of events and report them as configured. 

## Detection

A detection defines basics to build formatted messages for remote syslog server and holds reporters to analyze an SLH
event. If event matches conditions, it is getting reported as defined.

## Detection

| Attribute        | Type           | Info                          |
|------------------|:---------------|-------------------------------|
| `class_id`       | `<string>`     | Event Class ID                |
| `name`           | `<string>`     | Human readable message        |
| `severity`       | `<int>`        | Severity `0-10` (low to high) |
| `loglevel` &ast; | `<int>`        | Syslog Log Level (see below)  |
| `reporters`      | `[]<Reporter>` | Reporters to analyze events   |

&ast; Log levels:

```
LOG_EMERG   = 0
LOG_ALERT   = 1
LOG_CRIT    = 2
LOG_ERR     = 3
LOG_WARNING = 4
LOG_NOTICE  = 5
LOG_INFO    = 6
LOG_DEBUG   = 7
```

## Reporters

A reporter defines its type and passes in a configuration (`map[string]string`) for given type.

```json
{
  "type": "type_identifier",
  "config": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

Currently only a limited set of reporters is available.

### `type`

Check if event type matches `type` from config:

```json
{
  "type": "type",
  "config": {
    "type": "LOGIN"
  }
}
```

### `detail_exists`

Check if configured key `details` exists as key in event details:

```json
{
  "type": "detail_exists",
  "config": {
    "details": "my_event_property"
  }
}
```

### `detail_not_exists`

Check that configured key `details` does not exist as key in event details:

```json
{
  "type": "detail_not_exists",
  "config": {
    "details": "my_event_property"
  }
}
```
