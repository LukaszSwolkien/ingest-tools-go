# Ingest-Tools-Go
Sample Golang tools to ingest telemetry data into Splunk Observability suite

## Prerequisite
Use your organisation or create a new one using Splunk API and setup access tokens (see [here](https://github.com/LukaszSwolkien/ingest-tools)).

> **_NOTE:_**
You can create `.conf.yaml` file with url, token, protocol, transport etc. These parameters will be used as default. You can skip entering command line args that are defined in the `.conf.yaml`, that you do not want to overwrite, for example: 
```yaml
token: "ACCESS_TOKEN"
ingest: "trace"
protocol: "otlp"
transport: "grpc"
url: "ingest.REALM.signalfx.com:443"
```

# What samples are implemented

|Ingest    | Transport | Protocol        | ENDPOINT           | Sample  |
|----------|-----------|-----------------|--------------------|---------|
|trace     |   gRPC    | OTLP/trace/v1   | (port:443)         | &check; |
|trace     |   HTTP    | OTLP/trace/v1   | /v2/trace/otlp     | &cross; |
|trace     |   HTTP    | Zipkin JSON     | /v2/trace          | &check; |
|metrics   |   gRPC    | OTLP/metrics/v1 | _not implemented_  |  NA     |
|metrics   |   HTTP    | OTLP/metrics/v1 | /v2/datapoint/otlp | &cross; |
|metrics   |   HTTP    | SignalFx JSON   | /v2/datapoint      | &check; |
|log       |   gRPC    | OTLP/logs/v1    | _not implemented_  |  NA     |
|log       |   HTTP    | OTLP/logs/v1    | _not implemented_  |  NA     |
|log       |   HTTP    | Splunk HEC      | /v1/log            | &check; |
|events    |   gRPC    | OTLP/logs/v1    | _not implemented_  | NA      |
|events    |   HTTP    | OTLP/logs/v1    | v3/events          | &cross; |
|events    |   HTTP    | SignalFx        | v2/events          | &cross; |
|rum       |           | OTLP/logs/v1    |                    | &cross; |
|profiling |           | OTLP/logs/v1    |                    | &cross; |

## Protocols:

* [OTLP proto files](https://github.com/open-telemetry/opentelemetry-proto/tree/main/opentelemetry/proto) 
* [Zipkin JSON](https://zipkin.io/pages/data_model.html)
* [SignalFx JSON](https://dev.splunk.com/observability/reference/api/ingest_data/latest#endpoint-send-metrics)
* [Splunk HEC](https://docs.splunk.com/Documentation/Splunk/latest/Data/FormatEventsforHTTPEventCollector)

# Usage

ingest tool needs following parameters to run:
```bash
Usage:
    go run . -i=INGEST -p=PROTOCOL -t=TRANSPORT -url=URL -token=TOKEN [grpc-insecure=false]
Options:
    -i  The INGEST type (trace, metrics, logs, events, rum)
    -p  The request PROTOCOL (zipkin, otlp, sapm, thrift)
    -t  TRANSPORT (http, grpc)
    -token  Ingest access TOKEN
    -url    The URL to ingest endpoint
    -grpc-insecure  (optional) Set grpc-insecure=true to disable TLS
```

> **_NOTE_**: concatenate `https://ingest.REALM.signalfx.com` with `ENDPOINT` for ingest url (see above table) with HTTP transport. Use `ingest.REALM.signalfx.com` for gRPC calls.

## Examples:

* OTLP/gRPC trace sample:
```bash
go run . -i=trace -p=otlp -t=grpc -url=ingest.REALM.signalfx.com:443 -token=TOKEN
```

* Zipkin Json/HTTP trace sample:
```bash
go run . -i=trace -p=zipkin -t=http -url=https://ingest.lab0.signalfx.com/v2/trace -token=TOKEN
```

* SignalFx Json Datapoint/HTTP metrics sample:
```bash
go run . -i=metrics -p=sfx -t=http -url=https://ingest.REALM.signalfx.com/v2/datapoint -token=TOKEN
```

* OTLP/HTTP metrics sample:
```bash
go run . -i=metrics -p=otlp -t=http -url=https://ingest.REALM.signalfx.com/v2/datapoint/otlp -token=TOKEN
```

* Splunk HEC/HTTP log sample:
```bash
go run . -i=log -p=hec -t=http -url=https://ingest.REALM.signalfx.com/v1/logs -token=TOKEN
```

# Mock ingest services
You can use a mock server to consume samples instead of the actual endpoint, however you won't be able to see the sent data in Splunk Observability suite.

To run the mock trace-ingest service:

```bash
go run ./cmd/mock/trace-server
```