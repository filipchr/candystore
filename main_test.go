package main

import (
	"log"
	"os"
	"testing"
)

func TestParseData(t *testing.T) {
	data := [][]byte{
		[]byte("foo, Geisha, 100"),
		[]byte("foo, Dumle, 100"),
		[]byte("foo, Dumle, 50"),
		[]byte("foo, Lakrits, 200"),
	}
	tranactions := map[string]*Transaction{}

	for _, d := range data {
		parseData(d, tranactions)
	}

	if tranactions["foo"].Name != "foo" {
		t.Errorf("Not the same name")
	}
	if tranactions["foo"].Total != 450 {
		t.Errorf("Not the total")
	}

	if tranactions["foo"].CandyCount["Geisha"] != 1 ||
		tranactions["foo"].CandyCount["Dumle"] != 2 ||
		tranactions["foo"].CandyCount["Lakrits"] != 1 {
		t.Errorf("Not the Candy count")
	}

	if tranactions["foo"].TopSnack.Name != "Lakrits" {
		t.Errorf("Not the Candy count")
	}
}

func BenchmarkRun(b *testing.B) {
	b.StopTimer()
	fileData, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer fileData.Close()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		run(fileData)
	}

}
