package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	wg sync.WaitGroup
	NL byte = 10
)

type TopSnack struct {
	Name  string
	Count int
}
type Transaction struct {
	Name       string
	CandyCount map[string]int
	Total      int64
	TopSnack   *TopSnack
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name           string `json:"name"`
		FavouriteSnack string `json:"favouriteSnack"`
		TotalSnacks    int64  `json:"totalSnacks"`
	}{
		Name:           t.Name,
		FavouriteSnack: t.TopSnack.Name,
		TotalSnacks:    t.Total,
	})
}

func worker(wg *sync.WaitGroup, cs chan []byte, data []byte) {
	defer wg.Done()
	cs <- data
}

func monitorWorker(wg *sync.WaitGroup, cs chan []byte) {
	wg.Wait()
	close(cs)
}

func readFile(reader *bufio.Reader, records chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		data, err := reader.ReadBytes(NL)

		if err != nil || err == io.EOF {
			return
		}
		wg.Add(1)

		// remove the trailing new line
		data = data[:len(data)-1]

		go worker(wg, records, data)
	}
}

func parseData(data []byte, transactions map[string]*Transaction) {
	d := string(data)

	parts := strings.Split(d, ", ")

	name := parts[0]
	candy := parts[1]
	eaten, _ := strconv.ParseInt(parts[2], 10, 64)

	t, ok := transactions[name]

	if !ok {
		transactions[name] = &Transaction{
			Name:       name,
			CandyCount: map[string]int{candy: 1},
			Total:      eaten,
			TopSnack: &TopSnack{
				Name:  candy,
				Count: 1,
			},
		}
		return
	}

	t.Total = transactions[name].Total + eaten
	t.CandyCount[candy] = t.CandyCount[candy] + 1

	if candy != t.TopSnack.Name {
		t.TopSnack.Count = t.CandyCount[candy]
		t.TopSnack.Name = candy
	}
}

func run(file *os.File) []*Transaction {
	records := make(chan []byte)
	reader := bufio.NewReader(file)

	wg.Add(1)
	go readFile(reader, records, &wg)

	go monitorWorker(&wg, records)

	transactions := map[string]*Transaction{}
	for data := range records {
		parseData(data, transactions)
	}

	ledger := []*Transaction{}

	for _, t := range transactions {
		ledger = append(ledger, t)
	}

	sort.Slice(ledger, func(i, j int) bool {
		return ledger[i].Total > ledger[j].Total
	})

	return ledger
}

func main() {
	fileData, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer fileData.Close()

	data := run(fileData)

	result, _ := json.Marshal(data)
	fmt.Println(string(result))
}
