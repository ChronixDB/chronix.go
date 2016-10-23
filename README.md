[![Build Status](https://travis-ci.org/ChronixDB/chronix.go.svg?branch=master)](https://travis-ci.org/ChronixDB/chronix.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ChronixDB/chronix.go)](https://goreportcard.com/report/github.com/ChronixDB/chronix.go)
[![code-coverage](http://gocover.io/_badge/github.com/ChronixDB/chronix.go/chronix)](http://gocover.io/github.com/ChronixDB/chronix.go/chronix)
[![go-doc](https://godoc.org/github.com/ChronixDB/chronix.go/chronix?status.svg)](https://godoc.org/github.com/ChronixDB/chronix.go/chronix)
[![Apache License 2](http://img.shields.io/badge/license-ASF2-blue.svg)](https://github.com/ChronixDB/chronix.go/blob/master/LICENSE)

# The Chronix Go Client Library
This repository contains the Go client library for Chronix. It allows writing
time series data into Chronix and reading it back. While the write
implementation allows storing structured time series data, the read
implementation is still rudimentary and returns only an opaque byte slice
(usually containing JSON, but this depends on the `fl` query parameter).

For full details on usage, see the
[Go package documentation](https://godoc.org/github.com/ChronixDB/chronix.go/chronix).

# Example Usage

[This example](https://github.com/ChronixDB/chronix.go/blob/master/example)
stores several test time series in Chronix and reads them back.

## Importing

```go
import "github.com/ChronixDB/chronix.go/chronix"
```

## Creating a Chronix Client

```go
// Parse the Solr/Chronix URL.
u, err := url.Parse("http://<solr-url>/solr/chronix")
if err != nil {
	// Handle error.
}

// Create a Solr client.
solr := chronix.NewSolrClient(u, nil)

// Construct a Chronix client based on the Solr client.
c := chronix.New(solr)
```

## Writing Series Data

```go
// Construct a test time series with one data point.
series := []chronix.TimeSeries{
	{
		Metric: "testmetric",
		Attributes: map[string]string{
			"host": "testhost",
		},
		Points: []chronix.Point{
			{
				Timestamp: 1470784794,
				Value: 42.23,
			},
		},
	},
}

// Store the test series and commit within one second.
err := c.Store(series, false, time.Second)
if err != nil {
  // Handle error.
}
```

## Querying Series Data

```go
// Define the Chronix query parameters.
q := "metric:(testmetric) AND start:1470784794000 AND end:1470784794000"
fq := "join=host_s,metric"
fl := "dataAsJson"

// Execute the query.
resp, err := c.Query(q, fq, fl)
if err != nil {
  // Handle error.
}
```
