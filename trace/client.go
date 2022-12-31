package trace

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
)

func PostTraceSample(url string, secret string, contentType string, trace ZipkinTrace) {
	json_data, err := json.MarshalIndent(trace, "", "\t")
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	log.Println("* Ingest endpoint: " + url)
	log.Println("* Sending sample data:\n" + string(json_data))

	r, err := http.NewRequest("POST", url+"/v2/trace", bytes.NewBuffer(json_data))
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
