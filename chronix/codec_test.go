package chronix

import (
	"io/ioutil"
	"reflect"
	"testing"
	"encoding/base64"
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

// We've seen that different environment produce different prefixes. The test checks that both variants produce the
// same result.
func TestDecodeWithDifferentPrefixes(t *testing.T) {
	dataVariant1 := "H4sIAAAJbogA/0TULUgDYRzH8RNWvBkWDIsXVhTDJVl8gtGwJBNMYphNTAqCIAoKolNOxDHxbep8P3W+Tz0dyMKi0bC4YBLDJbEcv++lzz089397jsdu77Kix447bdHLoNH6qBiMszopmmk2zIjNOfYuiM4iny2L7goRVsWMR7A1cWyduBtivkCKouhvkm1LzG2TeEcMYbArzu9RTklM7lMZLB9Q5KGYLlPvkViH+WMxe0IXp2IL+mfixDm9XYgxnzZhA7qXogdDOHTFHGDqmpHAb5ipiBWYvBGnYBP23TI+2HHHJOEnTN+LRWg9iCOwDnsemTr8hdknsQqdqjgLW7D/mROCnS8cFvyC5lUswVgg9sIcLMAG/IPumzgMPfgBQ9j9zk8Al2AAf2CqZux4Iro3BmomYf0HAAD//zye/EFSBAAA"
	dataVariant2 := "H4sIAAAAAAAA/0TULUgDYRzH8RNWvBkWDIsXVhTDJVl8gtGwJBNMYphNTAqCIAoKolNOxDHxbep8P3W+Tz0dyMKi0bC4YBLDJbEcv++lzz089397jsdu77Kix447bdHLoNH6qBiMszopmmk2zIjNOfYuiM4iny2L7goRVsWMR7A1cWyduBtivkCKouhvkm1LzG2TeEcMYbArzu9RTklM7lMZLB9Q5KGYLlPvkViH+WMxe0IXp2IL+mfixDm9XYgxnzZhA7qXogdDOHTFHGDqmpHAb5ipiBWYvBGnYBP23TI+2HHHJOEnTN+LRWg9iCOwDnsemTr8hdknsQqdqjgLW7D/mROCnS8cFvyC5lUswVgg9sIcLMAG/IPumzgMPfgBQ9j9zk8Al2AAf2CqZux4Iro3BmomYf0HAAD//zye/EFSBAAA"


	dataDecoded1, err := base64.StdEncoding.DecodeString(dataVariant1)
	if err != nil {
		t.Fatal("Base64 decoding failed:", err)
	}
	dataDecoded2, err := base64.StdEncoding.DecodeString(dataVariant2)
	if err != nil {
		t.Fatal("Base64 decoding failed:", err)
	}

	points1, err := decode(dataDecoded1, 15, 114, 15, 114)
	if err != nil {
		t.Fatal("Decoding failed:", err)
	}

	if len(points1) == 0 {
		t.Error("Expecting more than one point")
	}

	points2, err := decode(dataDecoded2, 15, 114, 15, 114)
	if err != nil {
		t.Fatal("Decoding failed:", err)
	}

	if !reflect.DeepEqual(points1, points2) {
		t.Fatalf("Points are not equal. Want:\n\n%v\n\nGot:\n\n%v", points1, points2)
	}
}
