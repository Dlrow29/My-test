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
	"time"
)

type Downloader struct {
	concurrency int //并发下载数量
}

// NewDownloader 创建一个新的下载器实例
func NewDownloader(concurrency int) *Downloader {
	return &Downloader{
		concurrency: concurrency,
	}
}

func (d *Downloader) Download(url, outFile string) error {
	if outFile == "" {
		outFile = path.Base(url)
	} // 实现下载逻辑

	// 添加重试机制
	var resp *http.Response
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		fmt.Printf("尝试获取文件信息... (第%d次)\n", i+1)
		resp, err = http.Head(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		if err != nil {
			fmt.Printf("Head请求失败: %v\n", err)
		} else {
			fmt.Printf("服务器返回状态码: %d\n", resp.StatusCode)
		}

		if i < maxRetries-1 {
			fmt.Printf("等待2秒后重试...\n")
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		return fmt.Errorf("请求Head失败 (重试%d次后): %v", maxRetries, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("服务器返回错误状态: %s", resp.Status)
	}

	if resp.Header.Get("Accept-Ranges") == "bytes" {
		fmt.Printf("文件大小: %d 字节，支持断点续传\n", resp.ContentLength)
		return d.multiDownload(url, outFile, int(resp.ContentLength)) // 支持断点续传
	}

	fmt.Printf("文件大小: %d 字节，不支持断点续传\n", resp.ContentLength)
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
		start = end
	}
	wg.Wait()
	d.merge(outFile)
	fmt.Printf("下载完成，文件保存为: %s\n", outFile)
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
		_, err = io.CopyBuffer(destFile, partFile, make([]byte, 8*1024*1024))
		if err != nil {
			return err
		}
		partFile.Close()
		os.Remove(partFilename)
	}
	return nil
}

func (d *Downloader) singleDownload(url, outFile string) error {
	fmt.Println("不支持断点续传")
	return nil
}

func (d *Downloader) downloadPartial(url, outFile string, start, end, i int) {
	if start >= end {
		return
	}

	fmt.Printf("开始下载分片 %d: bytes=%d-%d\n", i, start, end)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("分片 %d 创建请求失败: %v", i, err)
		return
		req.Header.Set("Accept-Encoding", "identity")

	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; x86_64) Go-Downloader/1.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("分片 %d 请求失败: %v", i, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		log.Printf("分片 %d 状态码异常: %d", i, resp.StatusCode)
		return
	}

	partFile, err := os.OpenFile(d.getPartFilename(outFile, i), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("分片 %d 创建文件失败: %v", i, err)
		return
	}
	defer partFile.Close()

	buf := make([]byte, 8*1024*1024) // 保持1MB缓冲区提高下载速度
	_, err = io.CopyBuffer(partFile, resp.Body, buf)
	if err != nil {
		log.Printf("分片 %d 写入文件失败: %v", i, err)
		return
	}

	fmt.Printf("分片 %d 下载完成\n", i)
}

func (d *Downloader) getPartDir(outFile string) string {
	return strings.SplitN(outFile, ".", 2)[0]
}

func (d *Downloader) getPartFilename(outFile string, i int) string {
	partDir := d.getPartDir(outFile)
	return fmt.Sprintf("%s-%s-%d", partDir, outFile, i)
}
