package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/ChronixDB/chronix.go/chronix"
)

func main() {
	solrURL := flag.String("solr.url", "", "The URL to the Solr endpoint to use.")
	flag.Parse()

	if *solrURL == "" {
		log.Fatalln("Need to provide -solr.url flag")
	}
	u, err := url.Parse(*solrURL)
	if err != nil {
		log.Fatalln("Error parsing Solr URL:", err)
	}
	solr := chronix.NewSolrClient(u, nil)
	c := chronix.New(solr)

	series := make([]chronix.TimeSeries, 0, 10)
	for s := 0; s < 10; s++ {
		ts := chronix.TimeSeries{
			Metric: "testmetric",
			Attributes: map[string]string{
				"host": fmt.Sprintf("testhost_%d", s),
			},
		}

		tsStart := time.Now().UnixNano() / 1e6
		ts.Points = make([]chronix.Point, 0, 100)
		for i := 0; i < 100; i++ {
			ts.Points = append(ts.Points, chronix.Point{
				Timestamp: tsStart + int64(i+15),
				Value:     float64((s + i) * 100),
			})
		}

		series = append(series, ts)
	}

	if err = c.Store(series, true); err != nil {
		log.Fatalln("Error storing time series:", err)
	}
}
