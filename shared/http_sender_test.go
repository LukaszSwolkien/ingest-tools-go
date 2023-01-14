package shared

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/LukaszSwolkien/IngestTools/ut"
)

func TestClientSendJsonData(t *testing.T){
	type dummy struct {
		name string
	}

	svr := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			}))
	defer svr.Close()

	sc := SendData(svr.URL, "t0p_s3cr3t", "application/json", dummy{name: "dummy_data"})

	ut.AssertTrue(t, sc == 200)
}