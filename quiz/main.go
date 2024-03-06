package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Example: go run main.go --csv=<your csv file name> --timeout=<30> --shuffle=<true>

func main() {
	csvFilePath := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeout := flag.Int("timeout", 5, "timeout in seconds")
	doShuffle := flag.Bool("shuffle", false, "shuffle the questions")

	fmt.Printf("You have %d seconds to answer as many questions as you can get, press Enter to start:", *timeout)
	fmt.Scanln()

	flag.Parse()

	lines, err := readFile(*csvFilePath)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilePath))
	}

	if *doShuffle {
		rand.Shuffle(len(lines), func(i, j int) {
			lines[i], lines[j] = lines[j], lines[i]
		})
	}

	timer := time.NewTimer(time.Duration(*timeout) * time.Second)

	correctAnswers := 0
	answerChannel := make(chan string)

	for i, line := range lines {
		fmt.Printf("Problem #%d: %s = \n", i+1, line.q)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d.\n", correctAnswers, len(lines))
			return
		case answer := <-answerChannel:
			if answer == line.a {
				correctAnswers++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correctAnswers, len(lines))
}

type problem struct {
	q string
	a string
}

func readFile(filePath string) ([]problem, error) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	problems := make([]problem, len(lines))

	for i, line := range lines {
		problems[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return problems, nil
}

func exit(msg string) {
	fmt.Println(msg)
	fmt.Println("bebra")
	os.Exit(1)
}
