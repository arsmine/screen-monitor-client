# screen-monitor-client

## What's screen-monitor-client?
screen-monitor-client is a stand-alone program that checks given url of screens
whether are running or not.

## Installation
* You need Golang V1.9.7.
```bash
$ git clone https://github.com/arsmine/screen-monitor-client
$ go get -u github.com/arsmine/screen-monitor-client
$ cd screen-monitor-client
$ go build
```

### Or

You can get binaries for Linux64 and Windows64 from releases.

## Usage
`./screen-monitor-client --config example_config.json`

### config.json
* `interval`: Time interval(10s, 1m etc.) that how frequently you want to check screens.
* `urls`: Screen urls that you want to check.
* `triger`: Name of the script that runs when a screen is down.
* `trigerType`: Binary location that runs triger. 

```json
{
    "interval": "30s",
    "urls": [
        "http://myendpoint.com/api/screens"
    ],
    "triger": "sendmessage.sh",
    "trigerType": "/bin/bash"
}
```

## Command-Line Arguments

### required
* `--config <config.json>`

