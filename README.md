# Ingest-Tools-Go
Sample Golang tools to ingest telemetry data into Splunk Observability suite

## Usage examples:
Create your organisation using Splunk API (see [here](https://github.com/LukaszSwolkien/ingest-tools)), and setup access tokens.

To test Zipkin Json format:
```bash
go run ./main.go --token=my_token --ingest=https://ingest.REALM.signalfx.com --endpoint=v2/trace
```

To test OTLP format over gRPC transport:
```bash
go run main.go --token=my_token --endpoint=v2/trace/otlp --ingest=ingest.lab0.signalfx.com:443
```

# Setup project 
You can create `.conf.yaml` file with ingest url, token, endpoint etc defined instead of using command line args, for example:

```yaml
token: "YOUR_ACCESS_TOKEN"
ingest: "https://ingest.REALM.signalfx.com"
endpoint: "v2/trace" 
protocol: "zipkin"
```