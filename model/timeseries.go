package model

type TimeSeries struct {
	Metric     string
	Attributes map[string]string
	Points     []Point
}

type Point struct {
	Timestamp int64
	Value     float64
}
