package shared

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
)

// Converts data structure into json and sends to ingest 
func SendDataSample(url string, secret string, data interface{}) {
	contentType := "application/json"
	json_data, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	log.Println("Sending sample data:\n" + string(json_data))

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	r.Header.Add("Content-Type", contentType)
	r.Header.Add("X-SF-Token", secret)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}

	resp, err := client.Do(r)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	dr, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(dr))
}