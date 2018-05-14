package chronix

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"strings"
	"bufio"
	"os"
	"reflect"
	"path/filepath"
)

func createElasticMock(reference string, t *testing.T) *httptest.Server {
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

		referenceFile := filepath.Join("fixtures", "elastic", reference)

		body = normalizeDataBlocks(body[:])
		writeReferenceFile(body[:], referenceFile, t)


		want := readReference(referenceFile, t)
		got := parseReceived(body[:], t)

		if !reflect.DeepEqual(want, got) {
			t.Fatalf("Unexpected request body. Want:\n\n%v\n\nGot:\n\n%v", want, got)
		}
	}))
	return server
}

func writeReferenceFile(body []byte, referenceFile string, t *testing.T) {
	if *update {
		if err := ioutil.WriteFile(referenceFile, body[:], 0644); err != nil {
			t.Fatal("Writing the reference file failed:", err)
		}
	}
}

func parseReceived(body []byte, t *testing.T) []interface{}{
	var result []interface{}

	s := string(body[:])
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		var asJson interface{}
		err := json.Unmarshal(scanner.Bytes(), &asJson)
		if err != nil {
			t.Fatal(err)
		}
		result = append(result, asJson)
	}
	return result
}

func readReference(reference string, t *testing.T) []interface{} {
	var result []interface{}

	file, err := os.Open(reference)
	if err != nil {
		t.Fatal("Opening file failed", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var asJson interface{}
		err = json.Unmarshal(scanner.Bytes(), &asJson)
		if err != nil {
			t.Fatal(err)
		}
		result = append(result, asJson)
	}
	return result
}

func createElasticClient(server *httptest.Server, createStatistics bool) Client {
	elastic := NewElasticTestStorage(&server.URL)

	if createStatistics {
		return NewWithStatistics(elastic)
	} else {
		return New(elastic)
	}
}

func TestElasticUpdateEndToEnd(t *testing.T) {
	server := createElasticMock("reference.txt", t)
	defer server.Close()

	c := createElasticClient(server, false)

	series := genTimeSeries()

	if err := c.Store(series, true, time.Second); err != nil {
		t.Fatal("Error storing time series:", err)
	}
}

func TestElasticUpdateWithStatisticsEndToEnd(t *testing.T) {
	server := createElasticMock("referenceWithStatistics.txt", t)
	defer server.Close()

	c := createElasticClient(server, true)

	series := genTimeSeries()

	if err := c.Store(series, true, time.Second); err != nil {
		t.Fatal("Error storing time series:", err)
	}
}

func TestElasticQueryEndToEnd(t *testing.T) {
	//not yet implemented
}
