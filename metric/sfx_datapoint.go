package metric
import (
	"os"
	"math/rand"
)

var (
	counterValue int64 = 0
)

type Datapoint struct {
	Metric      string `json:"metric" binding:"required"`
	Value       int64 `json:"value" binding:"required"`
	Dimensions  *map[string]string `json:"dimensions,omitempty" bson:"dimensions,omitempty"`
	Timestamp	*uint64 `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

type GuageDatapoint struct {
	Guage []Datapoint `json:"guage" binding:"required"`
}

type CounterDatapoint struct {
	Counter []Datapoint `json:"counter" binding:"required"`
}

type CumulativeCounterDatapoint struct {
	CumulativeCounter []Datapoint `json:"cumulative_counter" binding:"required"`
}

func GenerateSfxGuageDatapointSample() GuageDatapoint {
	hostName, err := os.Hostname()
	if (err != nil){
		hostName="unknown"
	}
	g := []Datapoint{
		{
		Metric: "heartbeat", 
		Value: rand.Int63n(1000), 
		Dimensions: &map[string]string{
			"host": hostName,
			},
		},
	}
	return GuageDatapoint{Guage: g}
}

func counter() []Datapoint {
	counterValue++
	hostName, err := os.Hostname()
	if (err != nil){
		hostName="unknown"
	}
	return []Datapoint{
		{
		Metric: "number-of-requests", 
		Value: counterValue, 
		Dimensions: &map[string]string{
			"host": hostName,
			},
		},
	}
}

func GenerateSfxCounterDatapointSample() CounterDatapoint {
	return CounterDatapoint{Counter: counter()}
}

func GenerateSfxCumulativeCounterDatapointSample() CumulativeCounterDatapoint {
	return CumulativeCounterDatapoint{CumulativeCounter: counter()}
}
