# Ingest-Tools-Go
Sample Golang tools to ingest telemetry data into Splunk Observability suite

## Usage examples:
Create your organisation using Splunk API (see [here](https://github.com/LukaszSwolkien/ingest-tools)), and setup access tokens.
Remember to set all necessery secrets in the `.secrets.yaml` file

```bash
go run ./main.go
```

TODO: add OTLP examples for trace, metric and log data.

TODO: add missing examples for all ingest endpoints

# Setup project 
### Create `.secrets.yaml` file with tokens and endpoints defined, for example:

```yaml
splunk-ingest-token: "your_access_token"
splunk-ingest: "https://ingest.{realm}.signalfx.com"
```