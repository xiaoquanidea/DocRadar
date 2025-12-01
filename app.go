package main

import (
	"context"
	"doc-radar/scanner"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"sort"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	scanner  *scanner.Scanner
	exporter *scanner.Exporter
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		scanner:  scanner.NewScanner(),
		exporter: scanner.NewExporter(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 设置扫描进度回调
	a.scanner.SetProgressCallback(func(progress scanner.ScanProgress) {
		wailsRuntime.EventsEmit(ctx, "scan-progress", progress)
	})
}

// DriveInfo 驱动器信息
type DriveInfo struct {
	Path  string `json:"path"`
	Label string `json:"label"`
}

// GetDrives 获取系统驱动器列表
func (a *App) GetDrives() []DriveInfo {
	var drives []DriveInfo

	if goruntime.GOOS == "windows" {
		// Windows系统
		for _, drive := range "CDEFGHIJKLMNOPQRSTUVWXYZ" {
			drivePath := string(drive) + ":\\"
			if _, err := os.Stat(drivePath); err == nil {
				drives = append(drives, DriveInfo{
					Path:  drivePath,
					Label: string(drive) + ":",
				})
			}
		}
	} else {
		// macOS/Linux
		homeDir, err := os.UserHomeDir()
		if err == nil {
			drives = append(drives, DriveInfo{
				Path:  homeDir,
				Label: "用户目录 (" + filepath.Base(homeDir) + ")",
			})
		}
		drives = append(drives, DriveInfo{
			Path:  "/",
			Label: "根目录 (/)",
		})
		// 添加常用目录
		commonDirs := []struct {
			path  string
			label string
		}{
			{"/Volumes", "外部卷 (/Volumes)"},
			{"/Users", "用户 (/Users)"},
		}
		for _, dir := range commonDirs {
			if _, err := os.Stat(dir.path); err == nil {
				drives = append(drives, DriveInfo{
					Path:  dir.path,
					Label: dir.label,
				})
			}
		}
	}

	return drives
}

// SelectDirectory 打开目录选择对话框
func (a *App) SelectDirectory() (string, error) {
	return wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "选择要扫描的目录",
	})
}

// SelectExportDirectory 选择导出目录
func (a *App) SelectExportDirectory() (string, error) {
	return wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "选择导出目录",
	})
}

// ScanFiles 扫描文件
func (a *App) ScanFiles(options scanner.ScanOptions) (*scanner.ScanResult, error) {
	return a.scanner.Scan(options)
}

// ExportFiles 导出文件
func (a *App) ExportFiles(options scanner.ExportOptions) (*scanner.ExportResult, error) {
	return a.exporter.Export(options)
}

// GetExportProgress 获取导出进度
func (a *App) GetExportProgress() scanner.ExportProgress {
	return a.exporter.GetProgress()
}

// ExportAsZip 导出为压缩包
func (a *App) ExportAsZip(options scanner.ExportOptions) (*scanner.ExportResult, error) {
	return a.exporter.ExportAsZip(options)
}

// OpenFolder 打开文件所在文件夹
func (a *App) OpenFolder(filePath string) error {
	dir := filepath.Dir(filePath)

	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", "/select,", filePath)
	case "darwin":
		cmd = exec.Command("open", "-R", filePath)
	default: // linux
		cmd = exec.Command("xdg-open", dir)
	}

	return cmd.Start()
}

// FilterResult 过滤结果
type FilterResult struct {
	Files      []scanner.FileInfo `json:"files"`
	TotalCount int                `json:"totalCount"`
}

// FilterFiles 过滤文件列表
func (a *App) FilterFiles(files []scanner.FileInfo, filter FilterOptions) FilterResult {
	var result []scanner.FileInfo

	for _, file := range files {
		// 按文件类型过滤
		if len(filter.FileTypes) > 0 {
			found := false
			for _, t := range filter.FileTypes {
				if file.FileType == t {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// 按有效性过滤
		if filter.ValidOnly && !file.IsValid {
			continue
		}
		if filter.InvalidOnly && file.IsValid {
			continue
		}

		// 按文件名搜索
		if filter.SearchText != "" {
			if !containsIgnoreCase(file.Name, filter.SearchText) {
				continue
			}
		}

		// 按文件大小过滤
		if filter.MinSize > 0 && file.Size < filter.MinSize {
			continue
		}
		if filter.MaxSize > 0 && file.Size > filter.MaxSize {
			continue
		}

		result = append(result, file)
	}

	// 排序
	sortFiles(result, filter.SortBy, filter.SortDesc)

	return FilterResult{
		Files:      result,
		TotalCount: len(result),
	}
}

// FilterOptions 过滤选项
type FilterOptions struct {
	FileTypes   []string `json:"fileTypes"`
	ValidOnly   bool     `json:"validOnly"`
	InvalidOnly bool     `json:"invalidOnly"`
	SearchText  string   `json:"searchText"`
	MinSize     int64    `json:"minSize"`
	MaxSize     int64    `json:"maxSize"`
	SortBy      string   `json:"sortBy"`   // name, size, modTime, type
	SortDesc    bool     `json:"sortDesc"` // 是否降序
}

// containsIgnoreCase 忽略大小写的字符串包含检查
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(substr) == 0 ||
			findIgnoreCase(s, substr) >= 0)
}

func findIgnoreCase(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if equalIgnoreCase(s[i:i+len(substr)], substr) {
			return i
		}
	}
	return -1
}

func equalIgnoreCase(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

// sortFiles 对文件列表排序
func sortFiles(files []scanner.FileInfo, sortBy string, desc bool) {
	sort.Slice(files, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "name":
			less = files[i].Name < files[j].Name
		case "size":
			less = files[i].Size < files[j].Size
		case "modTime":
			less = files[i].ModTime.Before(files[j].ModTime)
		case "type":
			less = files[i].FileType < files[j].FileType
		default:
			less = files[i].Name < files[j].Name
		}
		if desc {
			return !less
		}
		return less
	})
}
