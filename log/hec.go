package log

import (
	"time"
	"os"
)

// Event metadata: https://docs.splunk.com/Documentation/Splunk/8.0.5/Data/FormateventsforHTTPEventCollector
type SplunkEvent struct {
	// Event Metadata
	Time 		int64				`json:"time"`
	Host		string 				`json:"host,omitempty"`
	Source 		string 				`json:"source,omitempty"`
	SourceType 	string 				`json:"sourcetype,omitempty"`
	Index	 	string 				`json:"index,omitempty"`
	Fields		map[string]string	`json:"fields,omitempty"`
	// Event Data
	EventData	EventData			`json:"event"`
}

type EventData struct {
	Message		string	`json:"message"`
	Severity	string	`json:"severity,omitempty"`
}

func GenerateLogSample() SplunkEvent {
	hostName, err := os.Hostname()
	if (err != nil){
		hostName="unknown"
	}
	return SplunkEvent {
		Time: time.Now().Unix(),
		Host: hostName,
		Source: "sample-log-generator",
		EventData: EventData{
			Message: "Something terrible happened",
			Severity: "ERROR",
		},
	}
}
