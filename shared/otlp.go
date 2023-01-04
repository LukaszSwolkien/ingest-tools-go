package shared

import (
	common "go.opentelemetry.io/proto/otlp/common/v1"				// OTLP commons (basic data representation)
	resource "go.opentelemetry.io/proto/otlp/resource/v1"			// OTLP resource (metadata)
)

func GetResource(serviceName string) *resource.Resource {
	return &resource.Resource{
			Attributes: []*common.KeyValue{
				{
					Key: "service.name",
					Value: &common.AnyValue {
						Value: &common.AnyValue_StringValue{
							StringValue: serviceName,
						},
					},
				},
			},
			DroppedAttributesCount: 0,
		}
}

func GetInstrumentationScope(name string) *common.InstrumentationScope{
	return &common.InstrumentationScope{
		Name:    name,
		Version: "1.0.0",
	}
}
