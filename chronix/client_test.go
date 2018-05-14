package chronix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
	"flag"
	"path/filepath"
	"strings"
)

var update = flag.Bool("update", false, "update reference json files")

// Helper function that creates a Solr mock
func createSolrMock(t *testing.T, reference string) *httptest.Server {

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/solr/chronix/update?commit=true&commitWithin=1000" {
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

		body = normalizeDataBlocks(body[:])

		var got interface{}
		if err = json.Unmarshal(body, &got); err != nil {
			t.Fatal("Error unmarshalling body:", err)
		}


		referenceFile := filepath.Join("fixtures", "solr", reference)
		writeReferenceJson(got, t, referenceFile)

		// Read and unmarshal fixture.
		f, err := ioutil.ReadFile(referenceFile)
		if err != nil {
			t.Fatal("Error reading fixture file:", err)
		}
		var want interface{}
		if err := json.Unmarshal(f, &want); err != nil {
			t.Fatal("Error unmarshalling fixture:", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Unexpected request body. Want:\n\n%v\n\nGot:\n\n%v", want.([]interface{}), got.([]interface{}))
		}
	}))
}

func writeReferenceJson(jsonStruct interface{}, t *testing.T, referenceFile string) {
	if *update {
		pretty, err := json.MarshalIndent(jsonStruct, "", "    ")
		if err != nil {
			t.Fatal("Error marshalling request:", err)
		}
		ioutil.WriteFile(referenceFile, pretty, 0644)
	}
}

func normalizeDataBlocks(data []byte) []byte {
	s := string(data[:])
	converted := []byte(strings.Replace(s, "H4sIAAAAAAAA/", "H4sIAAAJbogA/", -1))
	return converted
}

// Helper function that generates some test data
func genTimeSeries() []*TimeSeries {
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
	return series
}

func TestUpdateEndToEnd(t *testing.T) {
	// given:
	server := createSolrMock(t, "update.json")
	defer server.Close()
	series := genTimeSeries()
	c, err := createSolrClient(server, t, false)

	// expect:
	if err = c.Store(series, true, time.Second); err != nil {
		t.Fatal("Error storing time series:", err)
	}
}

func TestUpdateWithStatsEndToEnd(t *testing.T) {
	// given:
	server := createSolrMock(t, "updateWithStats.json")
	defer server.Close()
	series := genTimeSeries()
	c, err := createSolrClient(server, t, true)

	// expect:
	if err = c.Store(series, true, time.Second); err != nil {
		t.Fatal("Error storing time series:", err)
	}
}

func createSolrClient(server *httptest.Server, t *testing.T, withStatistics bool) (Client, error) {
	u, err := url.Parse(server.URL + "/solr/chronix")
	if err != nil {
		t.Fatal("Error parsing Solr URL:", err)
	}
	solr := NewSolrStorage(u, nil)
	if withStatistics {
		return NewWithStatistics(solr), err
	} else {
		return New(solr), err
	}
}


func TestQueryEndToEnd(t *testing.T) {
	q := "name:(testmetric) AND start:1471517965000 AND end:1471520557000"
	cj := "host_s,name"
	fl := "dataAsJson"

	resultJSON := "{}"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wantPath := "/solr/chronix/select"
		if r.URL.Path != wantPath {
			t.Fatalf("Unexpected path; want %s, got %s", wantPath, r.URL.Path)
		}
		if r.Method != "GET" {
			t.Fatal("Unexpected method; want GET, got", r.Method)
		}

		qs := r.URL.Query()
		wantParams := map[string]string{
			"q":  q,
			"cj": cj,
			"fl": fl,
			"wt": "json",
		}
		for k, v := range wantParams {
			if qs.Get(k) != v {
				t.Fatalf("Unexpected query param value; want %s, got %s", k, qs.Get(k))
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resultJSON))
	}))
	defer server.Close()

	u, err := url.Parse(server.URL + "/solr/chronix")
	if err != nil {
		t.Fatal("Error parsing Solr URL:", err)
	}
	solr := NewSolrStorage(u, nil)
	c := New(solr)

	res, err := c.Query(q, cj, fl)
	if err != nil {
		t.Fatal("Error querying:", err)
	}
	if string(res) != resultJSON {
		t.Fatalf("Unexpected result JSON; want %s, got %s", resultJSON, string(res))
	}
}
