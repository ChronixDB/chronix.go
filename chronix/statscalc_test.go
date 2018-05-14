package chronix

import "testing"

func TestCalculateStats(t *testing.T) {

	// given:
	series := TimeSeries{ Name: "Test", Type: "metric", Attributes: map[string]string{ "host": "node0"}}
	series.Points = []Point{{Timestamp: 123, Value: 1.5}, {Timestamp: 153, Value:1.9}, {Timestamp: 200, Value:0.2}}

	// when:
	stats, err := calculateStats(&series)

	// then:
	if err != nil {
		t.Fatal("Failed to calculate Stats", err)
	}

	if stats.count != 3 {
		t.Error("Expected count=3, got ", stats.count)
	}

	if stats.min != 0.2 {
		t.Error("Expected min=0.2, got ", stats.min)
	}

	if stats.max != 1.9 {
		t.Error("Expected max=1.9, got ", stats.max)
	}

	if stats.avg != 1.2 {
		t.Error("Expected avg=1.2, got ", stats.avg)
	}

	if stats.timespan != 77 {
		t.Error("Expected timespan=77, got ", stats.timespan)
	}
}

func TestCalculateStatsWithEmptyPoints(t *testing.T) {
	// given:
	series := TimeSeries{ Name: "Test", Type: "metric", Attributes: map[string]string{ "host": "node0"}}
	series.Points = []Point{}

	// when:
	_, err := calculateStats(&series)

	// then:
	if err == nil {
		t.Error("Expected exception, but was nil")
	}

	if err.Error() != "TimeSeries has no Points" {
		t.Error("Exception message is wrong: ", err.Error())
	}
}
