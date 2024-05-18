# DuckDNS updater
A simple DuckDNS updater meant to be run in a container.
It expects the config file to be mounted to `/data/config.json`. Logs will be written to `/data/logs`.


The config file should look like this:
```json
{
  "domains": ["yourdomain", "yourotherdomain"],
  "token": "064a0540-864c-4f0f-8bf5-23857452b0c1",
  "interval": 300,
  "log": true,
  "log_verbose": false,
}
```
_All fields are required._

#### `domains`
A list of domains to update.

#### `token`
Your DuckDNS token.

#### `interval`
The interval in seconds between updates/checks.

#### `log`
Whether to log IP changes to `/data/logs/log-{month}-{year}.log` or not.

#### `log_verbose`
_Has no effect if `log` is `false`._

Whether to log every check or only changes.
