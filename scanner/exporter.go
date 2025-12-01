package scanner

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ExportOptions 导出选项
type ExportOptions struct {
	DestPath      string     `json:"destPath"`      // 目标路径
	Files         []FileInfo `json:"files"`         // 要导出的文件列表
	KeepStructure bool       `json:"keepStructure"` // 是否保持目录结构
	Overwrite     bool       `json:"overwrite"`     // 是否覆盖已存在的文件
}

// ExportProgress 导出进度
type ExportProgress struct {
	Total     int     `json:"total"`
	Completed int     `json:"completed"`
	Failed    int     `json:"failed"`
	Current   string  `json:"current"`
	Percent   float64 `json:"percent"`
}

// ExportResult 导出结果
type ExportResult struct {
	Success      int      `json:"success"`
	Failed       int      `json:"failed"`
	FailedFiles  []string `json:"failedFiles"`
	SkippedFiles []string `json:"skippedFiles"`
}

// Exporter 文件导出器
type Exporter struct {
	progress     ExportProgress
	mu           sync.RWMutex
	progressChan chan ExportProgress
}

// NewExporter 创建新的导出器
func NewExporter() *Exporter {
	return &Exporter{
		progressChan: make(chan ExportProgress, 100),
	}
}

// Export 导出文件
func (e *Exporter) Export(options ExportOptions) (*ExportResult, error) {
	// 确保目标目录存在
	if err := os.MkdirAll(options.DestPath, 0755); err != nil {
		return nil, fmt.Errorf("无法创建目标目录: %w", err)
	}

	result := &ExportResult{
		FailedFiles:  make([]string, 0),
		SkippedFiles: make([]string, 0),
	}

	total := len(options.Files)
	var completed int64
	var failed int64

	// 使用工作池并发复制文件
	workerCount := 4
	fileChan := make(chan FileInfo, total)
	var wg sync.WaitGroup

	// 结果收集
	resultChan := make(chan struct {
		success bool
		file    string
		skipped bool
	}, total)

	// 启动工作协程
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				destPath := e.getDestPath(options, file)

				// 检查目标文件是否存在
				if !options.Overwrite {
					if _, err := os.Stat(destPath); err == nil {
						resultChan <- struct {
							success bool
							file    string
							skipped bool
						}{true, file.Path, true}
						atomic.AddInt64(&completed, 1)
						continue
					}
				}

				// 复制文件
				err := e.copyFile(file.Path, destPath)
				if err != nil {
					resultChan <- struct {
						success bool
						file    string
						skipped bool
					}{false, file.Path, false}
					atomic.AddInt64(&failed, 1)
				} else {
					resultChan <- struct {
						success bool
						file    string
						skipped bool
					}{true, file.Path, false}
				}
				atomic.AddInt64(&completed, 1)

				// 更新进度
				e.mu.Lock()
				e.progress = ExportProgress{
					Total:     total,
					Completed: int(atomic.LoadInt64(&completed)),
					Failed:    int(atomic.LoadInt64(&failed)),
					Current:   file.Name,
					Percent:   float64(atomic.LoadInt64(&completed)) / float64(total) * 100,
				}
				e.mu.Unlock()
			}
		}()
	}

	// 发送文件到工作池
	go func() {
		for _, file := range options.Files {
			fileChan <- file
		}
		close(fileChan)
	}()

	// 等待所有工作完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for r := range resultChan {
		if r.skipped {
			result.SkippedFiles = append(result.SkippedFiles, r.file)
			result.Success++
		} else if r.success {
			result.Success++
		} else {
			result.Failed++
			result.FailedFiles = append(result.FailedFiles, r.file)
		}
	}

	return result, nil
}

// getDestPath 获取目标路径
func (e *Exporter) getDestPath(options ExportOptions, file FileInfo) string {
	if options.KeepStructure {
		// 保持目录结构
		relPath := file.Path
		// 尝试获取相对路径
		for _, drive := range GetDrives() {
			if strings.HasPrefix(file.Path, drive) {
				relPath = strings.TrimPrefix(file.Path, drive)
				break
			}
		}
		return filepath.Join(options.DestPath, relPath)
	}

	// 不保持目录结构，直接放到目标目录
	return filepath.Join(options.DestPath, file.Name)
}

// copyFile 复制文件
func (e *Exporter) copyFile(src, dst string) error {
	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 同步到磁盘
	return dstFile.Sync()
}

// GetProgress 获取当前进度
func (e *Exporter) GetProgress() ExportProgress {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.progress
}

// ExportAsZip 导出为压缩包
func (e *Exporter) ExportAsZip(options ExportOptions) (*ExportResult, error) {
	result := &ExportResult{
		FailedFiles:  make([]string, 0),
		SkippedFiles: make([]string, 0),
	}

	// 生成压缩包文件名
	zipFileName := fmt.Sprintf("office-files-%s.zip", time.Now().Format("20060102-150405"))
	zipPath := filepath.Join(options.DestPath, zipFileName)

	// 确保目标目录存在
	if err := os.MkdirAll(options.DestPath, 0755); err != nil {
		return nil, fmt.Errorf("无法创建目标目录: %w", err)
	}

	// 创建压缩文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return nil, fmt.Errorf("无法创建压缩文件: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	total := len(options.Files)

	// 初始化进度
	e.mu.Lock()
	e.progress = ExportProgress{
		Total:     total,
		Completed: 0,
		Failed:    0,
		Current:   "",
		Percent:   0,
	}
	e.mu.Unlock()

	for i, file := range options.Files {
		// 更新进度
		e.mu.Lock()
		e.progress.Current = file.Name
		e.progress.Completed = i
		e.progress.Percent = float64(i) / float64(total) * 100
		e.mu.Unlock()

		// 确定压缩包内的路径
		var zipEntryPath string
		if options.KeepStructure {
			// 保持目录结构
			relPath := file.Path
			for _, drive := range GetDrives() {
				if strings.HasPrefix(file.Path, drive) {
					relPath = strings.TrimPrefix(file.Path, drive)
					break
				}
			}
			zipEntryPath = strings.TrimPrefix(relPath, string(filepath.Separator))
		} else {
			zipEntryPath = file.Name
		}

		// 添加文件到压缩包
		err := e.addFileToZip(zipWriter, file.Path, zipEntryPath)
		if err != nil {
			result.Failed++
			result.FailedFiles = append(result.FailedFiles, file.Path)
		} else {
			result.Success++
		}
	}

	// 完成进度
	e.mu.Lock()
	e.progress.Completed = total
	e.progress.Percent = 100
	e.mu.Unlock()

	return result, nil
}

// addFileToZip 添加文件到压缩包
func (e *Exporter) addFileToZip(zipWriter *zip.Writer, srcPath, zipEntryPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 获取文件信息
	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// 创建zip文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = zipEntryPath
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, srcFile)
	return err
}
