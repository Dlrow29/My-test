package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

type Config struct {
	URL         string //下载的URL
	outFile     string //保存的文件名
	concurrency int    //并发下载数量
}

func ParseCLI() Config {
	var cfg Config
	flag.StringVar(&cfg.URL, "u", "", "下载的URL")
	flag.StringVar(&cfg.outFile, "o", "", "保存到本地的文件名")
	flag.IntVar(&cfg.concurrency, "n", runtime.NumCPU(), "并发下载数量")
	flag.Usage = func() {
		_, err := fmt.Fprintf(os.Stderr, "Usage: %s -u <url> -o <output file> [-n <concurrency>]\n", os.Args[0])
		if err != nil {
			return
		}
		flag.PrintDefaults()
	}
	flag.Parse()

	if cfg.URL == "" || cfg.outFile == "" {
		log.Fatal("请提供下载的URL和保存的文件名")
	}
	if cfg.concurrency <= 0 {
		log.Fatal("并发下载数量必须大于0")
	}

	return cfg
}

func main() {
	cfg := ParseCLI()
	downloader := NewDownloader(cfg.concurrency)
	err := downloader.Download(cfg.URL, cfg.outFile)
	if err != nil {
		log.Fatal(err)
	}
}
