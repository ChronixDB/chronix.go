package codec

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/ChronixDB/chronix.go/model"
)

func gzipDecode(t *testing.T, r io.Reader) []byte {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		t.Fatal("Failed to create gzip reader: ", err)
	}
	defer gzipReader.Close()

	buf, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		t.Fatal("Failed to decompress data: ", err)
	}
	return buf
}

func buildTestPoints() []model.Point {
	points := make([]model.Point, 0, 100)
	for i := 0; i < 100; i++ {
		points = append(points, model.Point{
			Timestamp: int64(i + 15),
			Value:     float64(i * 100),
		})
	}
	return points
}

func TestEncode(t *testing.T) {
	points := buildTestPoints()

	buf, err := Encode(points)
	if err != nil {
		t.Fatal("Failed to encode points: ", err)
	}

	f, err := os.Open("fixtures/encoded.gz")
	if err != nil {
		t.Fatal("Failed to open test fixture: ", err)
	}
	defer f.Close()

	want := gzipDecode(t, f)
	got := gzipDecode(t, bytes.NewReader(buf))

	if !bytes.Equal(want, got) {
		t.Fatalf("wrong encoding; want:\n\n%v\n\ngot:\n\n%v", want, got)
	}
}

func TestDecode(t *testing.T) {
	points := buildTestPoints()

	buf, err := Encode(points)
	if err != nil {
		t.Fatal("Failed to encode points: ", err)
	}

	tsStart := points[0].Timestamp
	tsEnd := points[len(points)-1].Timestamp
	outPoints, err := Decode(buf, tsStart, tsEnd, tsStart, tsEnd)
	if err != nil {
		t.Fatal("Failed to decode points: ", err)
	}
	if !reflect.DeepEqual(points, outPoints) {
		t.Fatalf("Unexpected points, want:\n\n%v\n\ngot:\n\n%v", points, outPoints)
	}
}
