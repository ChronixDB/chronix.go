package chronix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"strings"
	"bufio"
)

func TestElasticUpdateEndToEnd(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.String() != "/_bulk" {
			t.Fatal("Unexpected URL:", r.URL.String())
		}

		if r.Method != "POST" {
			t.Fatal("Unexpected method:", r.Method)
		}

		// Read and unmarshal request body.
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal("Error reading request body:", err)
		}

		s := string(body[:])
		var lines []string

		scanner := bufio.NewScanner(strings.NewReader(s))
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		var got interface{}

		//Check that valid json is written
		for _, line := range lines {
			if err = json.Unmarshal([]byte(line), &got); err != nil {
				t.Fatal("Error unmarshalling body:", err)
			}
		}

	}))
	defer server.Close()

	solr := NewElasticTestStorage(&server.URL)
	c := New(solr)

	series := make([]*TimeSeries, 0, 10)
	for s := 0; s < 10; s++ {
		ts := &TimeSeries{
			Name: "testmetric",
			Type: "metric",
			Attributes: map[string]string{
				"host": fmt.Sprintf("testhost_%d", s),
			},
		}

		ts.Points = make([]Point, 0, 100)
		for i := 0; i < 100; i++ {
			ts.Points = append(ts.Points, Point{
				Timestamp: int64(i + 15),
				Value:     float64((s + i) * 100),
			})
		}

		series = append(series, ts)
	}

	if err := c.Store(series, true, time.Second); err != nil {
		t.Fatal("Error storing time series:", err)
	}
}

func TestElasticQueryEndToEnd(t *testing.T) {
	//not yet implemented
}
