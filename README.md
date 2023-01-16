# Ingest-Tools-Go
Sample Golang tools to ingest telemetry data into Splunk Observability suite

## Prerequisite
Use your organisation or create a new one using Splunk API and setup access tokens (see [here](https://github.com/LukaszSwolkien/ingest-tools)).

> **_NOTE:_**
You can create `.conf.yaml` file with url, token, protocol, transport etc. These parameters will be used as default. You can skip entering command line args that are defined in the `.conf.yaml`, that you do not want to overwrite, for example: 
```yaml
token: "ACCESS_TOKEN"
ingest: "trace"
data-format: "otlp"
transport: "grpc"
url: "ingest.REALM.signalfx.com:443"
```

# What samples are implemented

|Ingest type | Transport | Data-Format        | Endpoint               | Content-Type           | Sample  |
|------------|-----------|--------------------|------------------------|------------------------|---------|
|trace       |   gRPC    | OTLP/trace/v1      | (port:443)             | application/grpc       | &check; |
|trace       |   HTTP    | OTLP/trace/v1      | /v2/trace/otlp         | application/x-protobuf | &check; |
|trace       |   HTTP    | Zipkin JSON        | /v2/trace              | application/json       | &check; |
|trace       |   HTTP    | SAPM               | /v2/trace/sapm         | application/x-protobuf | &cross; |
|trace       |   HTTP    | SignalFx Trace     | /v2/trace/signalfxv1   | application/json       | &cross; |
|trace       |   HTTP    | JaegerThrift       | /v2/trace/jeagerthrift | application/x-thrift   | &cross; |
|metrics     |   HTTP    | OTLP/metrics/v1    | /v2/datapoint/otlp     | application/x-protobuf | &check; |
|metrics     |   HTTP    | SignalFx Datapoint | /v2/datapoint          | application/json       | &check; |
|log         |   HTTP    | Splunk HEC         | /v1/log                | application/json       | &check; |
|profiling   |   HTTP    | OTLP/logs/v1       | /v1/log                | application/json       | &cross; |
|log         |   HTTP    | Splunk HEC         | /services/collector    | application/json       | &cross; |
|events      |   HTTP    | OTLP/logs/v1       | /v3/events             | application/x-protobuf | &cross; |
|events      |   HTTP    | SignalFx Event     | /v2/events             | application/json       | &cross; |
|rum         |   HTTP    | Zipkin JSON        | /v1/rum                | application/json       | &cross; |
|rum         |   HTTP    | OTLP/logs/v1       | /v1/rumreplay          | application/x-protobuf | &cross; |
|rum         |   HTTP    | Zipkin JSON        | /v1/rumreplay          | application/json       | &cross; |


## Data formats:

* [OTLP proto files](https://github.com/open-telemetry/opentelemetry-proto/tree/main/opentelemetry/proto) 
* [SAPM (Splunk APM Protocol) ProtoBuf schema](https://github.com/signalfx/sapm-proto)
* [Zipkin JSON](https://zipkin.io/pages/data_model.html)
* [SignalFx Datapoint](https://dev.splunk.com/observability/reference/api/ingest_data/latest#endpoint-send-metrics)
* [SignalFx Event](https://dev.splunk.com/observability/reference/api/ingest_data/latest#endpoint-send-events)
* [SignalFx Trace](https://dev.splunk.com/observability/reference/api/ingest_data/latest#endpoint-sendtraces)
* [Splunk HEC](https://docs.splunk.com/Documentation/Splunk/latest/Data/FormatEventsforHTTPEventCollector)
* [Jaeger](https://www.jaegertracing.io/docs/1.41/apis/)

# Usage

ingest tool needs following parameters to run:
```bash
Usage:
    go run . -i=INGEST -f=FORMAT -t=TRANSPORT -url=URL -token=TOKEN [grpc-insecure=false]
Options:
    -i  The Ingest type (trace, metrics, logs, events, rum)
    -f  The request Data-Format (zipkin, otlp, sapm, thrift, sfx)
    -t  Transport (http, grpc)
    -token  Ingest access TOKEN
    -url    The URL to ingest endpoint
```

> **_NOTE_**: concatenate `https://ingest.REALM.signalfx.com` with `ENDPOINT` for ingest url (see above table) with HTTP transport. Use `ingest.REALM.signalfx.com` for gRPC calls.

## Examples:

* OTLP/gRPC trace sample:
```bash
go run . -i=trace -f=otlp -t=grpc -url=ingest.REALM.signalfx.com:443 -token=TOKEN
```

* Zipkin Json/HTTP trace sample:
```bash
go run . -i=trace -f=zipkin -t=http -url=https://ingest.lab0.signalfx.com/v2/trace -token=TOKEN
```

* SignalFx Json Datapoint/HTTP metrics sample:
```bash
go run . -i=metrics -f=sfx -t=http -url=https://ingest.REALM.signalfx.com/v2/datapoint -token=TOKEN
```

* OTLP/HTTP metrics sample:
```bash
go run . -i=metrics -f=otlp -t=http -url=https://ingest.REALM.signalfx.com/v2/datapoint/otlp -token=TOKEN
```

* Splunk HEC/HTTP log sample:
```bash
go run . -i=log -f=hec -t=http -url=https://ingest.REALM.signalfx.com/v1/logs -token=TOKEN
```

# Mock ingest services
You can use a mock server to consume samples instead of the actual endpoint, however you won't be able to see the sent data in Splunk Observability suite.

To run the mock trace-ingest service:

```bash
go run ./cmd/mock/trace-server
```

To connect with mock grpc service run on localhost use `grpc-insecure=true` flag to disable TLS

```
go run . -i=trace -f=otlp -t=grpc -url=localhost:8201 -grpc-insecure=true
```

# Testing
To run unit tests and check test coverage run below script
```bash
./test_cover.sh
```