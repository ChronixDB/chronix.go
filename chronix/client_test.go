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
)

func TestUpdateEndToEnd(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		var got interface{}
		if err = json.Unmarshal(body, &got); err != nil {
			t.Fatal("Error unmarshalling body:", err)
		}

		// Read and unmarshal fixture.
		f, err := ioutil.ReadFile("fixtures/update.json")
		if err != nil {
			t.Fatal("Error reading fixture file:", err)
		}
		var want interface{}
		if err := json.Unmarshal(f, &want); err != nil {
			t.Fatal("Error unmarshalling fixture:", err)
		}

		if !reflect.DeepEqual(got, want) {
			//Todo: Fix me.
			//t.Fatalf("Unexpected request body. Want:\n\n%v\n\nGot:\n\n%v", want, got)
		}
	}))
	defer server.Close()

	u, err := url.Parse(server.URL + "/solr/chronix")
	if err != nil {
		t.Fatal("Error parsing Solr URL:", err)
	}
	solr := NewSolrClient(u, nil)
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

	if err = c.Store(series, true, time.Second); err != nil {
		t.Fatal("Error storing time series:", err)
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
	solr := NewSolrClient(u, nil)
	c := New(solr)

	res, err := c.Query(q, cj, fl)
	if err != nil {
		t.Fatal("Error querying:", err)
	}
	if string(res) != resultJSON {
		t.Fatalf("Unexpected result JSON; want %s, got %s", resultJSON, string(res))
	}
}
