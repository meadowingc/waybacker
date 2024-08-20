# Wayback Machine Monitor

This project monitors a list of URLs using the Wayback Machine API. The URLs to monitor are specified in a YAML configuration file.

## Setup

1. Install Go: Follow the instructions at https://golang.org/doc/install to download and install Go.

2. Clone the repository: `git clone https://github.com/yourusername/wayback-machine-monitor.git`

3. Navigate to the project directory: `cd wayback-machine-monitor`

4. Build the project: `go build`

## Configuration

Grab your access key and secret key from https://archive.org/account/s3.php 

The configuration file is a YAML file named `config.yaml` with the following structure:

```yaml
access_key: your_access_key
secret_key: your_secret_key
sleep_days: 15
urls:
    - https://example.com
    - https://another-example.com
```