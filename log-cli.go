package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	logFilePaths []string
	totalFiles   int
	successFiles int
)

func init() {
	flag.Func("files", "Concurrent processing of files", func(values string) error {
		logFilePaths = strings.Split(values, ",")
		return nil
	})
	flag.Parse()
	totalFiles = len(logFilePaths)
}

var wg sync.WaitGroup

func ProcessFile(filePath string) {
	startTime := time.Now()

	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file %s: %v\n", filePath, err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "ERROR") {
			fmt.Println(line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Failed to read file %s: %v\n", filePath, err)
		return
	}

	log.Printf("Successfully processed file %s in %v\n", filePath, time.Since(startTime))
	successFiles++

}

func main() {
	start := time.Now()

	defer func() {
		duration := time.Since(start)
		log.Println("Total time taken:", duration)
		log.Printf("successfully rate %.2f%%\n", float64(successFiles)/float64(totalFiles)*100)
	}()
	for _, filePath := range logFilePaths {
		wg.Add(1)
		go ProcessFile(filePath)
	}

	wg.Wait()
	log.Println("All files processed successfully")

}
