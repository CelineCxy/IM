package main

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ImageTask struct {
	Name       string
	Data       []byte
	OutputPath string
	URL        string
	Err        error
}

func main() {
	// 三个阶段之间的传送带。
	images := make(chan ImageTask)
	compressedImages := make(chan ImageTask)
	uploadedResults := make(chan ImageTask)

	inputDir := `D:\IM\Images`
	outputDir := `D:\IM\Images\compressed`

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Println("创建输出目录失败：", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 阶段一：读取图片。
	go func() {
		defer close(images)

		entries, err := os.ReadDir(inputDir)
		if err != nil {
			fmt.Println("读取图片目录失败:", err)
			return
		}
		for _, entry := range entries {
			// 忽略文件夹和非PNG文件
			if entry.IsDir() {
				continue
			}
			if filepath.Ext(entry.Name()) != ".png" {
				continue
			}
			inputPath := filepath.Join(inputDir, entry.Name())
			data, err := os.ReadFile(inputPath)
			task := ImageTask{
				Name: entry.Name(),
				Data: data,
			}
			if err != nil {
				task.Err = fmt.Errorf("读取图片%s失败:%w", entry.Name(), err)
			}
			select {
			case images <- task:
				fmt.Println("读取图片:", inputPath)
			case <-ctx.Done():
				fmt.Println("读取阶段停止：", ctx.Err())
				return
			}
		}
	}()

	// 阶段二：压缩图片。
	var compressWG sync.WaitGroup
	for i := 0; i < 3; i++ {
		compressWG.Add(1)
		go func(i int) {
			defer compressWG.Done()
			for task := range images {
				fmt.Printf("压缩 worker %d 开始处理:%s\n", i, task.Name)
				if task.Err == nil {
					select {
					case <-ctx.Done():
						fmt.Printf("压缩 Worker %d 停止处理:%s\n", i, task.Name)
						return
					default:
					}
					task.OutputPath = filepath.Join(outputDir, "compressed-"+task.Name)
					err := compressPNG(task.Data, task.OutputPath)
					if err != nil {
						task.Err = fmt.Errorf("压缩图片 %s 失败: %w", task.Name, err)
					} else {
						fmt.Printf("压缩图片 %s 处理完成:%s\n", task.Name, task.OutputPath)
					}
					task.Data = nil
				}
				select {
				case compressedImages <- task:
				case <-ctx.Done():
					fmt.Printf("压缩 Worker %d 发送结果时停止\n", i)
					return
				}
			}
		}(i)
	}

	go func() {
		compressWG.Wait()
		close(compressedImages)
	}()

	// 阶段三：上传图片。
	var uploadWG sync.WaitGroup
	for i := 0; i < 2; i++ {
		uploadWG.Add(1)
		go func(workerID int) {
			defer uploadWG.Done()
			for task := range compressedImages {
				if task.Err != nil {
					fmt.Printf(
						"上传 Worker %d 跳过失败任务：%s\n", workerID, task.Name)
				} else {
					fmt.Printf("上传 Worker %d 开始上传：%s\n", workerID, task.OutputPath)

					if !wait(ctx, 500*time.Millisecond) {
						fmt.Printf("上传 Worker %d 停止\n", workerID)
						return
					}

					task.URL = "https://example.com/" + filepath.Base(task.OutputPath)
				}

				select {
				case uploadedResults <- task:
				case <-ctx.Done():
					fmt.Printf(
						"上传 Worker %d 发送结果时停止\n",
						workerID,
					)
					return
				}
			}
		}(i)
	}
	go func() {
		uploadWG.Wait()
		close(uploadedResults)
	}()

	// 阶段四：汇总结果。
	successCount := 0
	failedCount := 0
	for task := range uploadedResults {
		if task.Err != nil {
			failedCount++

			fmt.Printf("处理失败：name=%s error=%v\n", task.Name, task.Err)
			continue
		} else {
			successCount++

			fmt.Printf(
				"上传成功：name=%s url=%s\n",
				task.Name,
				task.URL,
			)
		}
	}
	if ctx.Err() != nil {
		fmt.Println("流水线提前结束：", ctx.Err())
		return
	}

	fmt.Printf("所有图片处理完成：成功 %d 张，失败 %d 张\n", successCount, failedCount)
}

func wait(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func compressPNG(data []byte, outputPath string) error {
	// 从内存中的 PNG 数据解码图片。
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("解码 PNG 失败: %w", err)
	}

	// 创建压缩后的输出文件。
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer outputFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	if err := encoder.Encode(outputFile, img); err != nil {
		return fmt.Errorf("编码 PNG 失败: %w", err)
	}

	return nil
}
