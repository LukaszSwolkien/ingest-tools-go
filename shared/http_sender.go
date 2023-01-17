package shared

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

// Converts data structure into json and sends to ingest
func SendJsonData(url string, secret string, contentType string, data interface{}) int {
	json_data, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Printf("Marshal: %v", err)
		return 400
	}
	log.Println("Sending sample data:\n" + string(json_data))

	body := bytes.NewBuffer(json_data)
	return PostHttpRequest(url, secret, contentType, body)
}

func PostHttpRequest(url string, secret string, contentType string, body io.Reader) int {
	r, err := http.NewRequest("POST", url, body)

	if err != nil {
		log.Fatalln(err)
	}
	r.Header.Add("Content-Type", contentType)
	r.Header.Add("X-SF-Token", secret)
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
	return resp.StatusCode
}
