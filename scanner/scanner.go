package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// FileInfo 文件信息结构
type FileInfo struct {
	Path          string    `json:"path"`
	Name          string    `json:"name"`
	Size          int64     `json:"size"`
	ModTime       time.Time `json:"modTime"`
	Extension     string    `json:"extension"`
	FileType      string    `json:"fileType"` // pdf, word, excel, ppt
	IsValid       bool      `json:"isValid"`
	InvalidReason string    `json:"invalidReason,omitempty"`
}

// ScanOptions 扫描选项
type ScanOptions struct {
	RootPath      string   `json:"rootPath"`
	IncludeTypes  []string `json:"includeTypes"`  // pdf, word, excel, ppt
	ExcludePaths  []string `json:"excludePaths"`  // 排除的路径
	ValidateFiles bool     `json:"validateFiles"` // 是否验证文件有效性
}

// ScanResult 扫描结果
type ScanResult struct {
	Files        []FileInfo `json:"files"`
	TotalCount   int        `json:"totalCount"`
	ValidCount   int        `json:"validCount"`
	InvalidCount int        `json:"invalidCount"`
	ScanTime     float64    `json:"scanTime"` // 扫描耗时（秒）
}

// ScanProgress 扫描进度
type ScanProgress struct {
	CurrentPath string `json:"currentPath"` // 当前扫描的路径
	ScannedDirs int    `json:"scannedDirs"` // 已扫描的目录数
	FoundFiles  int    `json:"foundFiles"`  // 已找到的文件数
	CurrentFile string `json:"currentFile"` // 当前处理的文件
	IsScanning  bool   `json:"isScanning"`  // 是否正在扫描
}

// ProgressCallback 进度回调函数类型
type ProgressCallback func(progress ScanProgress)

// Scanner 文件扫描器
type Scanner struct {
	options          ScanOptions
	mu               sync.Mutex
	progress         ScanProgress
	progressCallback ProgressCallback
}

// NewScanner 创建新的扫描器
func NewScanner() *Scanner {
	return &Scanner{}
}

// SetProgressCallback 设置进度回调
func (s *Scanner) SetProgressCallback(callback ProgressCallback) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.progressCallback = callback
}

// GetProgress 获取当前扫描进度
func (s *Scanner) GetProgress() ScanProgress {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.progress
}

// updateProgress 更新进度并通知
func (s *Scanner) updateProgress(progress ScanProgress) {
	s.mu.Lock()
	s.progress = progress
	callback := s.progressCallback
	s.mu.Unlock()

	if callback != nil {
		callback(progress)
	}
}

// 支持的文件扩展名映射
var extensionMap = map[string]string{
	// PDF
	".pdf": "pdf",
	// Word
	".doc":  "word",
	".docx": "word",
	".docm": "word",
	".dot":  "word",
	".dotx": "word",
	// Excel
	".xls":  "excel",
	".xlsx": "excel",
	".xlsm": "excel",
	".xlsb": "excel",
	".xlt":  "excel",
	".xltx": "excel",
	".csv":  "excel",
	// PowerPoint
	".ppt":  "ppt",
	".pptx": "ppt",
	".pptm": "ppt",
	".pot":  "ppt",
	".potx": "ppt",
	".pps":  "ppt",
	".ppsx": "ppt",
}

// 默认排除的路径关键词
var defaultExcludePaths = []string{
	// macOS
	"Library/Containers",
	"Library/Caches",
	"Library/Application Support",
	// 通用
	"node_modules",
	".git",
	".svn",
	// Windows 系统目录
	"$RECYCLE.BIN",
	"System Volume Information",
}

// Windows 系统目录前缀（需要精确匹配目录名）
var windowsSystemDirs = []string{
	"Windows",
	"Program Files",
	"Program Files (x86)",
	"ProgramData",
}

// Scan 执行文件扫描
func (s *Scanner) Scan(options ScanOptions) (*ScanResult, error) {
	startTime := time.Now()

	s.mu.Lock()
	s.options = options
	s.mu.Unlock()

	var files []FileInfo
	var mu sync.Mutex
	scannedDirs := 0
	foundFiles := 0
	lastUpdateTime := time.Now()

	// 初始化进度
	s.updateProgress(ScanProgress{
		CurrentPath: options.RootPath,
		ScannedDirs: 0,
		FoundFiles:  0,
		IsScanning:  true,
	})

	// 合并默认排除路径和用户指定的排除路径
	excludePaths := append(defaultExcludePaths, options.ExcludePaths...)

	// 检查路径是否应该被排除
	shouldExclude := func(path string, isDir bool) bool {
		// 统一路径分隔符为 /
		normalizedPath := strings.ReplaceAll(path, "\\", "/")

		// 检查通用排除路径
		for _, excludePath := range excludePaths {
			if strings.Contains(normalizedPath, excludePath) {
				return true
			}
		}

		// 检查 Windows 系统目录（只在根目录下排除）
		if isDir {
			dirName := filepath.Base(path)
			for _, sysDir := range windowsSystemDirs {
				if strings.EqualFold(dirName, sysDir) {
					// 检查是否是驱动器根目录下的系统目录
					parent := filepath.Dir(path)
					if len(parent) <= 3 { // 如 "C:\" 或 "C:"
						return true
					}
				}
			}
		}

		return false
	}

	err := filepath.WalkDir(options.RootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// 跳过无法访问的目录（权限问题、符号链接等）
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 检查是否应该排除此路径
		if shouldExclude(path, d.IsDir()) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 处理目录 - 更新进度
		if d.IsDir() {
			scannedDirs++
			// 每100ms更新一次进度，避免过于频繁
			if time.Since(lastUpdateTime) > 100*time.Millisecond {
				s.updateProgress(ScanProgress{
					CurrentPath: path,
					ScannedDirs: scannedDirs,
					FoundFiles:  foundFiles,
					IsScanning:  true,
				})
				lastUpdateTime = time.Now()
			}
			return nil
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(path))
		fileType, ok := extensionMap[ext]
		if !ok {
			return nil
		}

		// 检查是否在包含的类型中
		if len(options.IncludeTypes) > 0 {
			found := false
			for _, t := range options.IncludeTypes {
				if t == fileType {
					found = true
					break
				}
			}
			if !found {
				return nil
			}
		}

		// 获取文件信息
		info, err := d.Info()
		if err != nil {
			return nil
		}

		fileInfo := FileInfo{
			Path:      path,
			Name:      d.Name(),
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			Extension: ext,
			FileType:  fileType,
			IsValid:   true,
		}

		// 验证文件有效性
		if options.ValidateFiles {
			valid, reason := ValidateFile(path, fileType)
			fileInfo.IsValid = valid
			fileInfo.InvalidReason = reason
		}

		mu.Lock()
		files = append(files, fileInfo)
		foundFiles++
		mu.Unlock()

		// 更新进度（找到文件时）
		if time.Since(lastUpdateTime) > 100*time.Millisecond {
			s.updateProgress(ScanProgress{
				CurrentPath: filepath.Dir(path),
				ScannedDirs: scannedDirs,
				FoundFiles:  foundFiles,
				CurrentFile: d.Name(),
				IsScanning:  true,
			})
			lastUpdateTime = time.Now()
		}

		return nil
	})

	// 扫描完成，更新进度
	s.updateProgress(ScanProgress{
		CurrentPath: options.RootPath,
		ScannedDirs: scannedDirs,
		FoundFiles:  foundFiles,
		IsScanning:  false,
	})

	if err != nil {
		return nil, err
	}

	// 统计结果
	validCount := 0
	invalidCount := 0
	for _, f := range files {
		if f.IsValid {
			validCount++
		} else {
			invalidCount++
		}
	}

	return &ScanResult{
		Files:        files,
		TotalCount:   len(files),
		ValidCount:   validCount,
		InvalidCount: invalidCount,
		ScanTime:     time.Since(startTime).Seconds(),
	}, nil
}

// GetDrives 获取系统所有驱动器（Windows）或根目录（macOS/Linux）
func GetDrives() []string {
	var drives []string

	// macOS/Linux
	if _, err := os.Stat("/"); err == nil {
		// 获取用户主目录
		homeDir, err := os.UserHomeDir()
		if err == nil {
			drives = append(drives, homeDir)
		}
		drives = append(drives, "/")
	}

	// Windows - 检查常见驱动器盘符
	for _, drive := range "CDEFGHIJKLMNOPQRSTUVWXYZ" {
		drivePath := string(drive) + ":\\"
		if _, err := os.Stat(drivePath); err == nil {
			drives = append(drives, drivePath)
		}
	}

	return drives
}
