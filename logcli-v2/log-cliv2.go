package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	logFilePaths []string
	totalFiles   int
	successFiles int64
)

func init() {
	flag.Func("files", "comma-separated list of log files to scan", func(v string) error {
		logFilePaths = strings.Split(v, ",")
		return nil
	})
	flag.Parse()
	totalFiles = len(logFilePaths)
}

func processFile(ctx context.Context, filePath string, results chan<- string) error {
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("open %s: %v", filePath, err)
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if strings.Contains(scanner.Text(), "ERROR") {
			results <- scanner.Text()
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("read %s: %v", filePath, err)
		return err
	}
	atomic.AddInt64(&successFiles, 1)
	return nil
}

func main() {
	start := time.Now()
	const maxWorkers = 4
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobs := make(chan string, len(logFilePaths))
	results := make(chan string, 1000)
	sem := make(chan struct{}, maxWorkers)

	go func() {
		defer close(jobs)
		for _, p := range logFilePaths {
			select {
			case jobs <- p:
			case <-ctx.Done():
				return
			}
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				select {
				case sem <- struct{}{}:
					_ = processFile(ctx, path, results)
					<-sem
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for line := range results {
		fmt.Println(line)
	}

	log.Printf("总耗时: %v, 成功率: %.2f%%",
		time.Since(start),
		float64(atomic.LoadInt64(&successFiles))/float64(totalFiles)*100)
}
