package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

type DownloadResult struct {
	Bytes int64
	Err   error
}

// 下载一个 URL，读取完整响应体。
func download(client *http.Client, url string) DownloadResult {
	resp, err := client.Get(url)
	if err != nil {
		return DownloadResult{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return DownloadResult{
			Err: fmt.Errorf("HTTP 状态异常：%s", resp.Status),
		}
	}

	n, err := io.Copy(io.Discard, resp.Body)
	return DownloadResult{Bytes: n, Err: err}
}

// 串行：上一个下载完成后，才开始下一个。
func downloadSerial(client *http.Client,urls []string,) ([]DownloadResult, time.Duration) {
	results := make([]DownloadResult, len(urls))
	start := time.Now()

	for i, url := range urls {
		results[i] = download(client, url)
	}

	return results, time.Since(start)
}

// 并发：每个 URL 由一个 goroutine 下载。
func downloadConcurrent(client *http.Client,urls []string,) ([]DownloadResult, time.Duration) {
	results := make([]DownloadResult, len(urls))
	var wg sync.WaitGroup
	start := time.Now()

	for i, url := range urls {
		// 必须在启动 goroutine 前增加计数。
		wg.Add(1)

		go func(index int, downloadURL string) {
			defer wg.Done()

			// 每个 goroutine 写自己的元素，不修改切片长度。
			results[index] = download(client, downloadURL)
		}(i, url)
	}

	// 等待全部下载结束，再统计耗时和读取结果。
	wg.Wait()

	return results, time.Since(start)
}

func printResults(name string, results []DownloadResult, elapsed time.Duration) {
	success := 0
	var totalBytes int64

	for i, result := range results {
		if result.Err != nil {
			fmt.Printf("%s：第 %d 个 URL 下载失败：%v\n",
				name, i+1, result.Err)
			continue
		}

		success++
		totalBytes += result.Bytes
	}

	fmt.Printf("%s：成功 %d/%d，总字节数 %d，耗时 %v\n",
		name, success, len(results), totalBytes, elapsed)
}


func main() {
	// 本地模拟下载服务：每个请求耗时约 200ms。
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			timer := time.NewTimer(200 * time.Millisecond)
			defer timer.Stop()

			select {
			case <-timer.C:
				fmt.Fprintf(w, "这是 %s 的文件内容\n", r.URL.Path)
			case <-r.Context().Done():
				return
			}
		},
	))
	defer server.Close()

	urls := make([]string, 10)
	for i := range urls {
		urls[i] = fmt.Sprintf("%s/files/%d", server.URL, i+1)
	}
	// 也可以把 urls 替换为 10 个实际的 HTTP/HTTPS 地址。

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	defer client.CloseIdleConnections()

	serialResults, serialTime := downloadSerial(client, urls)
	printResults("串行下载", serialResults, serialTime)

	concurrentResults, concurrentTime := downloadConcurrent(client, urls)
	printResults("并发下载", concurrentResults, concurrentTime)

	fmt.Printf("串行耗时 / 并发耗时：%.2f\n",
		float64(serialTime)/float64(concurrentTime))
}