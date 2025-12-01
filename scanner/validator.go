package scanner

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

// 文件魔数（Magic Numbers）定义
var (
	// PDF 文件头
	pdfMagic = []byte{0x25, 0x50, 0x44, 0x46} // %PDF

	// Office 2007+ (OOXML) - ZIP 格式
	zipMagic = []byte{0x50, 0x4B, 0x03, 0x04} // PK..

	// Office 97-2003 (OLE2 Compound Document)
	oleMagic = []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}
)

// ValidateFile 验证文件是否有效
func ValidateFile(path string, fileType string) (bool, string) {
	file, err := os.Open(path)
	if err != nil {
		return false, "无法打开文件: " + err.Error()
	}
	defer file.Close()

	// 获取文件大小
	stat, err := file.Stat()
	if err != nil {
		return false, "无法获取文件信息: " + err.Error()
	}

	// 空文件检查
	if stat.Size() == 0 {
		return false, "文件为空"
	}

	// 读取文件头部用于验证
	header := make([]byte, 8)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		return false, "无法读取文件头: " + err.Error()
	}
	if n < 4 {
		return false, "文件太小，无法验证"
	}

	switch fileType {
	case "pdf":
		return validatePDF(file, header, stat.Size())
	case "word", "excel", "ppt":
		return validateOffice(file, header, stat.Size(), fileType)
	}

	return true, ""
}

// validatePDF 验证PDF文件
func validatePDF(file *os.File, header []byte, size int64) (bool, string) {
	// 检查PDF魔数
	if !bytes.HasPrefix(header, pdfMagic) {
		return false, "不是有效的PDF文件（文件头不匹配）"
	}

	// 检查PDF结尾标记 %%EOF
	// 读取文件末尾
	if size < 10 {
		return false, "PDF文件太小"
	}

	tailSize := int64(1024)
	if size < tailSize {
		tailSize = size
	}

	_, err := file.Seek(-tailSize, io.SeekEnd)
	if err != nil {
		return false, "无法读取文件尾部"
	}

	tail := make([]byte, tailSize)
	_, err = file.Read(tail)
	if err != nil {
		return false, "无法读取文件尾部"
	}

	// 检查是否包含 %%EOF 标记
	if !bytes.Contains(tail, []byte("%%EOF")) {
		return false, "PDF文件可能已损坏（缺少EOF标记）"
	}

	// 检查是否包含基本的PDF结构
	file.Seek(0, io.SeekStart)
	content := make([]byte, min(size, 4096))
	file.Read(content)

	// 检查是否有对象定义
	if !bytes.Contains(content, []byte("obj")) {
		return false, "PDF文件结构异常"
	}

	return true, ""
}

// validateOffice 验证Office文件
func validateOffice(file *os.File, header []byte, size int64, fileType string) (bool, string) {
	// 检查是否是OOXML格式（.docx, .xlsx, .pptx等）
	if bytes.HasPrefix(header, zipMagic) {
		return validateOOXML(file, size, fileType)
	}

	// 检查是否是OLE2格式（.doc, .xls, .ppt等）
	if bytes.HasPrefix(header, oleMagic) {
		return validateOLE2(file, size, fileType)
	}

	// CSV文件特殊处理
	if fileType == "excel" {
		file.Seek(0, io.SeekStart)
		content := make([]byte, min(size, 4096))
		n, _ := file.Read(content)
		if n > 0 && isValidCSV(content[:n]) {
			return true, ""
		}
	}

	return false, "不是有效的Office文件（文件头不匹配）"
}

// validateOOXML 验证OOXML格式文件
func validateOOXML(file *os.File, size int64, fileType string) (bool, string) {
	// OOXML是ZIP格式，检查ZIP结构
	file.Seek(0, io.SeekStart)

	// 读取更多内容来验证ZIP结构
	content := make([]byte, min(size, 8192))
	n, err := file.Read(content)
	if err != nil && err != io.EOF {
		return false, "无法读取文件内容"
	}

	// 检查是否包含Office特有的内容类型文件
	contentTypes := []byte("[Content_Types].xml")
	if !bytes.Contains(content[:n], contentTypes) {
		// 可能是加密的Office文件或损坏的文件
		// 检查是否有其他ZIP条目
		if !bytes.Contains(content[:n], []byte("PK")) {
			return false, "Office文件结构异常"
		}
	}

	// 根据文件类型检查特定内容
	switch fileType {
	case "word":
		if bytes.Contains(content[:n], []byte("word/")) ||
			bytes.Contains(content[:n], []byte("document.xml")) {
			return true, ""
		}
	case "excel":
		if bytes.Contains(content[:n], []byte("xl/")) ||
			bytes.Contains(content[:n], []byte("workbook.xml")) ||
			bytes.Contains(content[:n], []byte("sheet")) {
			return true, ""
		}
	case "ppt":
		if bytes.Contains(content[:n], []byte("ppt/")) ||
			bytes.Contains(content[:n], []byte("presentation.xml")) ||
			bytes.Contains(content[:n], []byte("slide")) {
			return true, ""
		}
	}

	// 如果有Content_Types.xml，认为是有效的Office文件
	if bytes.Contains(content[:n], contentTypes) {
		return true, ""
	}

	return false, "Office文件内容异常"
}

// validateOLE2 验证OLE2格式文件
func validateOLE2(file *os.File, size int64, fileType string) (bool, string) {
	// OLE2格式验证
	file.Seek(0, io.SeekStart)

	// 读取OLE2头部
	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil || n < 512 {
		return false, "无法读取OLE2头部"
	}

	// 验证OLE2签名
	if !bytes.HasPrefix(header, oleMagic) {
		return false, "OLE2签名不匹配"
	}

	// 检查扇区大小（通常是512或4096）
	sectorSize := binary.LittleEndian.Uint16(header[30:32])
	if sectorSize != 9 && sectorSize != 12 { // 2^9=512, 2^12=4096
		return false, "OLE2扇区大小异常"
	}

	// 检查文件大小是否合理
	if size < 1536 { // 至少需要3个扇区
		return false, "OLE2文件太小"
	}

	return true, ""
}

// isValidCSV 检查是否是有效的CSV文件
func isValidCSV(content []byte) bool {
	// 检查是否包含可打印字符
	printableCount := 0
	totalCount := len(content)

	for _, b := range content {
		// 可打印ASCII字符、换行、回车、制表符
		if (b >= 32 && b <= 126) || b == '\n' || b == '\r' || b == '\t' {
			printableCount++
		}
		// UTF-8多字节字符
		if b >= 0x80 {
			printableCount++
		}
	}

	// 如果超过90%是可打印字符，认为是有效的CSV
	return float64(printableCount)/float64(totalCount) > 0.9
}

// min 返回两个int64中的较小值
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
