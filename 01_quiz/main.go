package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Reads a path to file and returns parsed contents of CSV file.
// For any errors, it will return nil
func parseCSV(fileName string) [][]string {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Failed to open file (%s) due to following error: %+v\n", fileName, err)
		return nil
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Printf("Failed to read csv file due to error: %+v\n", err)
		return nil
	}

	return records
}

func main() {
	fileName := flag.String("input", "problems.csv", "CSV file of problems")
	timeLimit := flag.Int("timer", 30, "Time limit for quiz")
	flag.Parse()

	score := 0
	records := parseCSV(*fileName)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	answerChannel := make(chan int)
	for _, data := range records {
		fmt.Printf("%s: ", data[0])

		go func() {
			var response int
			fmt.Scan(&response)
			answerChannel <- response
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou answered %d out of %d questions\n", score, len(records))
			return
		case response := <-answerChannel:
			answer, err := strconv.Atoi(data[1])
			if err != nil {
				fmt.Printf("Failed to convert %s to integer\n", data[1])
				continue
			}
			if response == answer {
				score++
			}
		}
	}
	fmt.Printf("You answered %d out of %d questions\n", score, len(records))
}
