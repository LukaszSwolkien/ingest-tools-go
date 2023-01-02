# Ingest-Tools-Go
Sample Golang tools to ingest telemetry data into Splunk Observability suite

## Usage examples:
Create your organisation using Splunk API (see [here](https://github.com/LukaszSwolkien/ingest-tools)), and setup access tokens.

```bash
go run ./main.go -token=my_token -ingest=https://ingest.REALM.signalfx.com -endpoint=v2/trace
```

TODO: add missing examples for all ingest endpoints

# Setup project 
You can create `.conf.yaml` file with ingest url, token, endpoint etc defined instead of using command line args, for example:

```yaml
token: "YOUR_ACCESS_TOKEN"
ingest: "https://ingest.REALM.signalfx.com"
endpoint: "v2/trace" 
protocol: "zipkin"
```