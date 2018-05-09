package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/ChronixDB/chronix.go/chronix"
)

func buildSeries() []*chronix.TimeSeries {
	series := make([]*chronix.TimeSeries, 0, 10)
	for s := 0; s < 10; s++ {
		ts := &chronix.TimeSeries{
			Name: "testmetric",
			Type: "metric",
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
	return series
}

func setupSolr(storageUrl *string) chronix.Client {
	u, err := url.Parse(*storageUrl)
	if err != nil {
		log.Fatalln("Error parsing Solr URL:", err)
	}
	solrStorage := chronix.NewSolrStorage(u, nil)
	return chronix.New(solrStorage)
}

func setupElastic(storageUrl *string, withIndex *bool, deleteIndexIfExist *bool, sniffElasticNodes *bool) chronix.Client {
	elasticStorage := chronix.NewElasticStorage(storageUrl, withIndex, deleteIndexIfExist, sniffElasticNodes)
	return chronix.New(elasticStorage)
}

func main() {
	storageUrl := flag.String("url", "", "The URL to the Solr endpoint to use.")
	kind := flag.String("kind", "", "Kind: solr or elastic")
	esWithIndex := flag.Bool("es.withIndex", true, "Creates an index if it do not exists")
	esDeleteIndexIfExists := flag.Bool("es.deleteIndexIfExists", false, "Deletes the index if one exists (only in use with es.withIndex)")
	esSniff := flag.Bool("es.sniffNodes", false, "Should the elastic client sniff for nodes (only in use with kind 'elastic')")
	flag.Parse()

	if *storageUrl == "" {
		log.Fatalln("Need to provide -url flag")
	}

	var client chronix.Client

	if *kind == "solr" {
		client = setupSolr(storageUrl)
	} else if *kind == "elastic" {
		client = setupElastic(storageUrl, esWithIndex, esDeleteIndexIfExists, esSniff)
	} else {
		log.Fatalln("Need to provide valid -kind flag")
	}

	log.Println("Storing time series...")
	series := buildSeries()
	err := client.Store(series, true, 0)
	if err != nil {
		log.Fatalln("Error storing time series:", err)
	}
	log.Println("Done storing.")

	log.Println("Querying time series...")
	q := "name:(testmetric) AND start:1471517965000 AND end:NOW"
	cj := "host_s,name"
	fl := "dataAsJson"
	resp, err := client.Query(q, cj, fl)
	if err != nil {
		log.Fatalln("Error querying time series:", err)
	}
	log.Println("Raw query output:", string(resp))
}
