# DuckDNS updater
A simple DuckDNS updater meant to be run in a container.
It assumes full control over the `/data` directory. Make sure to mount it to a persistent volume.
It requires a `config.json` file in the `/data` directory. The file should look like this:
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
Whether to log IP changes to `/data/log-{month}-{year}.log` or not.

#### `log_verbose`
_Has no effect if `log` is `false`._

Whether to log every check or only changes.
