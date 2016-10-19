package chronix

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func buildTestPoints() []Point {
	points := make([]Point, 0, 100)
	for i := 0; i < 100; i++ {
		points = append(points, Point{
			Timestamp: int64(i + 15),
			Value:     float64(i * 100),
		})
	}
	return points
}

func TestDecode(t *testing.T) {
	points := buildTestPoints()

	encoded, err := ioutil.ReadFile("fixtures/encoded.gz")
	if err != nil {
		t.Fatal("Failed to read test fixture: ", err)
	}

	tsStart := points[0].Timestamp
	tsEnd := points[len(points)-1].Timestamp
	outPoints, err := decode(encoded, tsStart, tsEnd, tsStart, tsEnd)
	if err != nil {
		t.Fatal("Failed to decode points: ", err)
	}
	if !reflect.DeepEqual(points, outPoints) {
		t.Fatalf("Unexpected points, want:\n\n%v\n\ngot:\n\n%v", points, outPoints)
	}
}

func TestEncode(t *testing.T) {
	points := buildTestPoints()

	buf, err := encode(points, 0)
	if err != nil {
		t.Fatal("Failed to encode points: ", err)
	}

	tsStart := points[0].Timestamp
	tsEnd := points[len(points)-1].Timestamp
	outPoints, err := decode(buf, tsStart, tsEnd, tsStart, tsEnd)
	if err != nil {
		t.Fatal("Failed to decode points: ", err)
	}
	if !reflect.DeepEqual(points, outPoints) {
		t.Fatalf("Unexpected points, want:\n\n%v\n\ngot:\n\n%v", points, outPoints)
	}
}
