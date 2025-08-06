package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type Downloader struct {
	concurrency int //并发下载数量
}

func (d *Downloader) Download(url, outFile string) error {
	if outFile == "" {
		outFile = path.Base(url)
	} // 实现下载逻辑
	resp, err := http.Head(url)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK && resp.Header.Get("Accept-Ranges") == "bytes" {
		return d.multiDownload(url, outFile, int(resp.ContentLength)) // 支持断点续传
	}
	return d.singleDownload(url, outFile) // 不支持断点续传
}

func (d *Downloader) multiDownload(url, outFile string, contentLength int) error {
	partSize := contentLength / d.concurrency

	partDir := d.getPartDir(outFile)
	os.Mkdir(partDir, 0777)
	defer os.RemoveAll(partDir)

	var wg sync.WaitGroup
	wg.Add(d.concurrency)
	start := 0
	for i := 0; i < d.concurrency; i++ {
		go func(i, start int) {
			defer wg.Done()
			end := start + partSize
			if d.concurrency-1 == i {
				end = contentLength
			}
			d.downloadPartial(url, outFile, start, end, i)
		}(i, start)
		start += partSize + 1
	}
	wg.Wait()
	d.merge(outFile)
	return nil
}

func (d *Downloader) merge(outFile string) error {
	destFile, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()
	for i := 0; i < d.concurrency; i++ {
		partFilename := d.getPartFilename(outFile, i)
		partFile, err := os.Open(partFilename)
		if err != nil {
			return err
		}
		defer partFile.Close()
		io.Copy(destFile, partFile)
		os.Remove(partFilename)
	}
	return nil
}

func (d *Downloader) singleDownload(url, outFile string) error {
	return nil
}

func (d *Downloader) downloadPartial(url, outFile string, start, end, i int) {
	if start >= end {
		return
	}
	req, err := http.NewRequest("Get", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	partFile, err := os.OpenFile(d.getPartFilename(outFile, i), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer partFile.Close()
	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(partFile, resp.Body, buf)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal(err)
	}
}

func (d *Downloader) getPartDir(outFile string) string {
	return strings.SplitN(outFile, ".", 2)[0]
}

func (d *Downloader) getPartFilename(outFile string, i int) string {
	partDir := d.getPartDir(outFile)
	return fmt.Sprintf("%s-%s-%d", partDir, outFile, i)
}
