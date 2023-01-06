# Ingest-Tools-Go
Sample Golang tools to ingest telemetry data into Splunk Observability suite

## Prerequisite
Use your organisation or create a new one using Splunk API and setup access tokens (see [here](https://github.com/LukaszSwolkien/ingest-tools)).

> **_NOTE:_**
You can create `.conf.yaml` file with url, token, protocol, transport etc. These parameters will be used as default. You can skip entering command line args that are defined in the `.conf.yaml`, that you do not want to overwrite, for example: 
```yaml
token: "ACCESS_TOKEN"
url: "https://ingest.REALM.signalfx.com:443"
protocol: "otlp"
```

# Samples

|Ingest    | Transport | Protocol        | ENDPOINT           | Sample  |
|----------|-----------|-----------------|--------------------|---------|
|trace     |   gRPC    | OTLP/trace/v1   | :443               | &check; |
|trace     |   HTTP    | OTLP/trace/v1   | /v2/trace/otlp     | &cross; |
|trace     |   HTTP    | Zipkin JSON     | /v2/trace          | &check; |
|metrics   |   gRPC    | OTLP/metrics/v1 | _not implemented_  |  NA     |
|metrics   |   HTTP    | OTLP/metrics/v1 | /v2/datapoint/otlp | &cross; |
|metrics   |   HTTP    | SignalFx JSON   | /v2/datapoint      | &check; |
|log       |   gRPC    | OTLP/logs/v1    | _not implemented_  |  NA     |
|log       |   HTTP    | OTLP/logs/v1    | _not implemented_  |  NA     |
|log       |   HTTP    | Splunk HEC      | /v1/log            | &check; |
|event     |   gRPC    | OTLP/logs/v1    | _not implemented_  | NA      |
|event     |   HTTP    | OTLP/logs/v1    | v3/events          | &cross; |
|event     |   HTTP    | SignalFx        | v2/events          | &cross; |
|rum       |           | OTLP/logs/v1    |                    | &cross; |
|profiling |           | OTLP/logs/v1    |                    | &cross; |

## Protocols:

* [OTLP proto files](https://github.com/open-telemetry/opentelemetry-proto/tree/main/opentelemetry/proto) 
* [Zipkin JSON](https://zipkin.io/pages/data_model.html)
* [SignalFx JSON](https://dev.splunk.com/observability/reference/api/ingest_data/latest#endpoint-send-metrics)
* [Splunk HEC](https://docs.splunk.com/Documentation/Splunk/latest/Data/FormatEventsforHTTPEventCollector)

## Examples:

> **_NOTE_**: concatenate `https://ingest.REALM.signalfx.com` with `ENDPOINT` for ingest url (see above table)

* OTLP/gRPC trace sample:
```bash
go run . --ingest=trace --protocol=otlp --transport=grpc --token=TOKEN --url=ingest.REALM.signalfx.com:443
```

* Zipkin Json/HTTP trace sample:
```bash
go run . --ingest=trace --protocol=zipkin --transport=http --token=TOKEN --url=https://ingest.REALM.signalfx.com/v2/trace
```

* SignalFx Json Datapoint/HTTP metrics sample:
```bash
go run . --ingest=metrics --protocol=sfx --transpoer=http --token=TOKEN --url=https://ingest.REALM.signalfx.com/v2/datapoint
```

* OTLP/HTTP metrics sample:
```bash
go run . --ingest=metrics --protocol=otlp --transpoer=http --token=TOKEN --url=https://ingest.REALM.signalfx.com/v2/datapoint/otlp
```

* Splunk HEC/HTTP log sample:
```
go run . --ingest=log --protocol=hec --transport=http --token=TOKEN --url=https://ingest.REALM.signalfx.com/v1/logs
```