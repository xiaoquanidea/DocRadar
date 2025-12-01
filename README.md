# 🔍 DocRadar - 文档雷达

快速扫描和管理 Office 文档的桌面应用，帮助你找到散落在各处的 PDF、Word、Excel、PPT 文件。

[![Release](https://img.shields.io/github/v/release/xiaoquanidea/DocRadar)](https://github.com/xiaoquanidea/DocRadar/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Windows-lightgrey.svg)](#-快速开始)
[![Go](https://img.shields.io/badge/Go-1.23-00ADD8.svg)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3-4FC08D.svg)](https://vuejs.org/)
[![Build](https://github.com/xiaoquanidea/DocRadar/actions/workflows/build.yml/badge.svg)](https://github.com/xiaoquanidea/DocRadar/actions)

## ✨ 核心功能

- 🔍 **智能扫描** - 快速扫描指定目录，找出所有 Office 文档
- 📁 **多格式支持** - PDF、Word (.doc/.docx)、Excel (.xls/.xlsx)、PowerPoint (.ppt/.pptx)
- ✅ **文件验证** - 自动检测损坏或无效的文件
- 📊 **可视化界面** - 直观的图形界面，分页浏览，支持搜索过滤
- 📦 **批量导出** - 支持导出到文件夹或打包为 ZIP 压缩包
- 🖱️ **快捷操作** - 点击路径直接打开文件所在文件夹
- 🌍 **跨平台** - 支持 macOS 和 Windows

## 📸 应用截图

*（待添加截图）*

## 🚀 快速开始

### 📥 下载使用

👉 前往 [Releases](https://github.com/xiaoquanidea/DocRadar/releases) 页面下载最新版本

#### macOS

1. 下载对应版本：
   - Apple Silicon (M1/M2/M3/M4): `DocRadar-macOS-arm64.zip`
   - Intel Mac: `DocRadar-macOS-amd64.zip`

2. 如遇到"无法打开"或"无法验证开发者"警告：
   ```bash
   # 移除隔离属性
   xattr -cr ~/Downloads/DocRadar.app
   ```
   或者：右键点击应用 → 打开 → 仍要打开

#### Windows

1. 下载 `DocRadar-Windows-amd64.zip`
2. 解压后双击 `DocRadar.exe` 运行
3. 如遇到 Windows Defender 警告，点击"更多信息" → "仍要运行"

### 使用步骤

1️⃣ **选择扫描路径** - 选择驱动器或自定义目录

2️⃣ **选择文件类型** - 勾选要扫描的文件类型（PDF、Word、Excel、PPT）

3️⃣ **开始扫描** - 点击"开始扫描"按钮，等待扫描完成

4️⃣ **查看结果** - 浏览扫描结果，支持搜索和过滤

5️⃣ **导出文件** - 选择文件后导出到文件夹或压缩包

## 🛠️ 开发指南

### 环境要求

- Go 1.21+
- Node.js 18+
- Wails v2.11.0

### 安装依赖

```bash
# 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 安装前端依赖
cd frontend && npm install
```

### 开发模式

```bash
# 启动开发服务器（支持热重载）
wails dev
```

### 构建应用

#### 一键构建（推荐）

```bash
./build.sh
```

生成：
- `dist/doc-radar-x.x.x-macos-arm64.zip` - macOS Apple Silicon 版本
- `dist/doc-radar-x.x.x-macos-amd64.zip` - macOS Intel 版本
- `dist/doc-radar-x.x.x-windows-amd64.zip` - Windows 64位 版本
- `dist/doc-radar-x.x.x-windows-386.zip` - Windows 32位 版本

#### 手动构建

```bash
# macOS ARM64 版本（Apple Silicon）
wails build -platform darwin/arm64

# macOS AMD64 版本（Intel）
wails build -platform darwin/amd64

# Windows AMD64 版本
wails build -platform windows/amd64
```

> **注意**: 在 macOS 上构建 Windows 版本需要先安装 MinGW-w64：
> ```bash
> brew install mingw-w64
> ```

## 🏗️ 项目架构

```
doc-radar/
├── scanner/                    # Go 后端 - 扫描模块
│   ├── scanner.go             # 文件扫描引擎
│   ├── validator.go           # 文件有效性验证
│   └── exporter.go            # 文件导出功能
├── frontend/                   # Vue 3 前端
│   ├── src/
│   │   ├── App.vue            # 主应用
│   │   └── main.ts            # 入口文件
│   └── package.json
├── build/                      # 构建资源
│   ├── appicon.png            # 应用图标
│   ├── darwin/                # macOS 资源
│   └── windows/               # Windows 资源
├── app.go                      # Wails 应用绑定
├── main.go                     # 应用入口
├── build.sh                    # 构建脚本
├── wails.json                  # Wails 配置
└── README.md                   # 本文档
```

## 🔧 技术栈

| 层级 | 技术 | 版本 |
|------|------|------|
| 后端 | Go | 1.21+ |
| 前端框架 | Vue | 3.x |
| 语言 | TypeScript | 5.x |
| UI 组件 | Element Plus | 2.x |
| 构建工具 | Vite | 5.x |
| 桌面框架 | Wails | v2.11.0 |

## 🎯 支持的文件格式

### PDF
- `.pdf`

### Word
- `.doc`, `.docx`, `.docm`, `.dot`, `.dotx`

### Excel
- `.xls`, `.xlsx`, `.xlsm`, `.xlsb`, `.xlt`, `.xltx`, `.csv`

### PowerPoint
- `.ppt`, `.pptx`, `.pptm`, `.pot`, `.potx`, `.pps`, `.ppsx`

## ⚠️ 注意事项

### 安全提示

- **导出会复制文件** - 导出操作会将文件复制到目标位置，不会删除原文件
- **验证功能** - 启用"验证文件有效性"可以过滤损坏的文件

### 扫描说明

- 扫描会自动跳过系统目录（Windows、Program Files 等）
- 扫描会跳过 node_modules、.git 等开发目录
- 大目录扫描可能需要较长时间

## 🐛 常见问题

**Q: 为什么扫描很慢？**

A: 扫描速度取决于目录大小和文件数量。建议选择具体的项目目录而不是整个驱动器。

**Q: 为什么有些文件显示"无效"？**

A: 启用文件验证后，程序会检查文件的魔数（Magic Number）来判断文件是否损坏。

**Q: macOS 提示"无法验证开发者"怎么办？**

A: 运行 `xattr -cr /path/to/DocRadar.app` 或右键打开。

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [Wails](https://wails.io/) - 优秀的 Go + Web 桌面应用框架
- [Element Plus](https://element-plus.org/) - 强大的 Vue 3 UI 组件库
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架

---

Made with ❤️ by [xiaoquanidea](https://github.com/xiaoquanidea)
