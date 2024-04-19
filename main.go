package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var wg sync.WaitGroup

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
		data, _, err := reader.ReadLine()

		if err != nil || err == io.EOF {
			return
		}
		wg.Add(1)
		go worker(wg, records, data)
	}

}

func parseData(data []byte) {
	fmt.Println(data)
}

func run(file *os.File) {

	records := make(chan []byte)
	reader := bufio.NewReader(file)

	wg.Add(1)
	go readFile(reader, records, &wg)

	go monitorWorker(&wg, records)

	counter := 1
	for data := range records {
		fmt.Printf("%d ", counter)
		parseData(data)

		counter++
	}
}

func main() {
	data, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer data.Close()

	run(data)
}
