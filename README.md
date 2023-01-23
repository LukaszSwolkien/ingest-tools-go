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
|trace       |   HTTP    | SignalFx Trace     | /v2/trace/signalfxv1   | application/json       | &check; |
|trace       |   HTTP    | SAPM               | /v2/trace/sapm         | application/x-protobuf | &check; |
|trace       |   HTTP    | JaegerThrift       | /v2/trace/jeagerthrift | application/x-thrift   | &check; |
|metrics     |   HTTP    | OTLP/metrics/v1    | /v2/datapoint/otlp     | application/x-protobuf | &check; |
|metrics     |   HTTP    | SignalFx Datapoint | /v2/datapoint          | application/json       | &check; |
|log         |   HTTP    | Splunk HEC         | /v1/log                | application/json       | &check; |
|profiling   |   HTTP    | OTLP/logs/v1       | /v1/log                | application/json       |         |
|log         |   HTTP    | Splunk HEC         | /services/collector    | application/json       |         |
|events      |   HTTP    | OTLP/logs/v1       | /v3/events             | application/x-protobuf |         |
|events      |   HTTP    | SignalFx Event     | /v2/events             | application/json       |         |
|rum         |   HTTP    | Zipkin JSON        | /v1/rum                | application/json       |         |
|rum         |   HTTP    | OTLP/logs/v1       | /v1/rumreplay          | application/x-protobuf |         |
|rum         |   HTTP    | Zipkin JSON        | /v1/rumreplay          | application/json       |         |


## Data formats:

* [OTLP proto files](https://github.com/open-telemetry/opentelemetry-proto/tree/main/opentelemetry/proto) 
* [SAPM (Splunk APM Protocol) ProtoBuf schema](https://github.com/signalfx/sapm-proto)
* [Zipkin JSON](https://zipkin.io/pages/data_model.html)
* [SignalFx Datapoint](https://dev.splunk.com/observability/reference/api/ingest_data/latest#endpoint-send-metrics)
* [SignalFx Event](https://dev.splunk.com/observability/reference/api/ingest_data/latest#endpoint-send-events)
* [SignalFx Trace](https://dev.splunk.com/observability/reference/api/ingest_data/latest#endpoint-sendtraces)
* [Splunk HEC](https://docs.splunk.com/Documentation/Splunk/latest/Data/FormatEventsforHTTPEventCollector)
* [Jaeger](https://www.jaegertracing.io/docs/1.41/apis/)

# Main Usage

ingest tool needs following parameters to run:
```bash
Usage:
    go run . -i=INGEST -f=FORMAT -t=TRANSPORT -url=URL -token=TOKEN [grpc-insecure=false]
Options:
    -i  The Ingest type (trace, metrics, logs, events, rum)
    -f  The request Data-Format (zipkin, otlp, sapm, jaegerthrift, sfx)
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
# Sampler

Instead of sending samples directly to the endpoint you can serialise payload to the file and use curl for http requests. To do that you need to generate payload for a given data format. 

sampler tool needs following parameters to run:
```bash
Usage:
    go run ./cmd/sampler -i=INGEST -f=FORMAT -file=FILENAME
Options:
    -i      The Ingest type (trace, metrics, logs, events, rum)
    -f      The request Data-Format (zipkin, otlp, sapm, jaegerthrift, sfx,...)
    -file   Output file name for payload data (default: payload.data)
```

* Example for trace Jaeger Thrift data format over http:

```bash
 go run ./cmd/sampler -i="trace" -f="jaegerthrift" -file="payload.data"
```
than use curl to post http request with binary data:
```
curl -X POST https://ingest.REALM.signalfx.com/v2/trace -H "Content-Type: application/x-thrift" -H "X-SF-Token: ACCESS_TOKEN" --data-binary @payload.data -i
```

* Example for metrics OTLP protobuf data format over http:

```bash
 go run ./cmd/sampler -i="trace" -f="otlp" -file="metrics_otlp.bin"
```

```
curl -X POST https://ingest.lab0.signalfx.com/v2/datapoint/otlp -H "Content-Type: application/x-protobuf" -H "X-SF-Token: ACCESS_TOKEN" --data-binary @metrics_otlp.bin -i
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