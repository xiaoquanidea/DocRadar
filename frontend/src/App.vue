<template>
  <div class="app-container">
    <!-- 顶部工具栏 -->
    <el-header class="header">
      <h1 class="title">
        <el-icon><Document /></el-icon>
        DocRadar - 文档雷达
      </h1>
    </el-header>

    <el-container class="main-container">
      <!-- 左侧控制面板 -->
      <el-aside width="320px" class="sidebar">
        <el-card class="control-card">
          <template #header>
            <div class="card-header">
              <el-icon><Search /></el-icon>
              <span>扫描设置</span>
            </div>
          </template>

          <!-- 扫描路径选择 -->
          <div class="form-section">
            <label class="form-label">扫描路径</label>
            <el-radio-group v-model="scanMode" class="scan-mode-group">
              <el-radio label="drive">选择驱动器</el-radio>
              <el-radio label="custom">自定义路径</el-radio>
            </el-radio-group>

            <div class="path-selector">
              <el-select
                v-if="scanMode === 'drive'"
                v-model="selectedDrive"
                placeholder="选择驱动器"
                class="full-width"
              >
                <el-option
                  v-for="drive in drives"
                  :key="drive.path"
                  :label="drive.label"
                  :value="drive.path"
                />
              </el-select>

              <div v-else class="path-input">
                <el-input
                  v-model="customPath"
                  placeholder="输入或选择路径"
                  readonly
                />
                <el-button @click="selectDirectory" type="primary">
                  <el-icon><FolderOpened /></el-icon>
                </el-button>
              </div>
            </div>
          </div>

          <!-- 文件类型选择 -->
          <div class="form-section">
            <label class="form-label">文件类型</label>
            <el-checkbox-group v-model="selectedTypes">
              <el-checkbox label="pdf">
                <el-tag type="danger" size="small">PDF</el-tag>
              </el-checkbox>
              <el-checkbox label="word">
                <el-tag type="primary" size="small">Word</el-tag>
              </el-checkbox>
              <el-checkbox label="excel">
                <el-tag type="success" size="small">Excel</el-tag>
              </el-checkbox>
              <el-checkbox label="ppt">
                <el-tag type="warning" size="small">PPT</el-tag>
              </el-checkbox>
            </el-checkbox-group>
          </div>

          <!-- 验证选项 -->
          <div class="form-section">
            <el-checkbox v-model="validateFiles">
              验证文件有效性（过滤损坏文件）
            </el-checkbox>
          </div>

          <!-- 扫描按钮 -->
          <el-button
            type="primary"
            size="large"
            class="full-width scan-btn"
            :loading="scanning"
            @click="startScan"
          >
            <el-icon v-if="!scanning"><Search /></el-icon>
            {{ scanning ? '扫描中...' : '开始扫描' }}
          </el-button>

          <!-- 扫描进度 -->
          <div v-if="scanning" class="scan-progress">
            <div class="progress-stats">
              <span>已扫描目录: {{ scanProgress.scannedDirs }}</span>
              <span>已找到文件: {{ scanProgress.foundFiles }}</span>
            </div>
            <div class="progress-bar-container">
              <div class="progress-bar-animated"></div>
            </div>
            <div class="current-path" :title="scanProgress.currentPath">
              <el-icon><Folder /></el-icon>
              <span>{{ truncatePath(scanProgress.currentPath) }}</span>
            </div>
          </div>
        </el-card>

        <!-- 过滤器和统计信息 -->
        <el-card class="control-card" v-if="scanResult">
          <template #header>
            <div class="card-header">
              <el-icon><Filter /></el-icon>
              <span>过滤与统计</span>
            </div>
          </template>

          <!-- 搜索框 -->
          <div class="form-section">
            <el-input
              v-model="filterText"
              placeholder="搜索文件名..."
              clearable
              @input="applyFilter"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>

          <!-- 有效性过滤 -->
          <div class="form-section">
            <el-radio-group v-model="validityFilter" @change="applyFilter" size="small">
              <el-radio-button label="all">全部</el-radio-button>
              <el-radio-button label="valid">有效</el-radio-button>
              <el-radio-button label="invalid">无效</el-radio-button>
            </el-radio-group>
          </div>

          <!-- 统计信息 -->
          <div class="stats-grid">
            <div class="stat-item">
              <span class="stat-value">{{ scanResult.totalCount }}</span>
              <span class="stat-label">总文件</span>
            </div>
            <div class="stat-item">
              <span class="stat-value valid">{{ scanResult.validCount }}</span>
              <span class="stat-label">有效</span>
            </div>
            <div class="stat-item">
              <span class="stat-value invalid">{{ scanResult.invalidCount }}</span>
              <span class="stat-label">无效</span>
            </div>
            <div class="stat-item">
              <span class="stat-value selected">{{ selectedFiles.length }}</span>
              <span class="stat-label">已选</span>
            </div>
          </div>

          <!-- 扫描耗时 -->
          <div class="scan-time">
            <el-icon><Timer /></el-icon>
            <span>扫描耗时: {{ scanResult.scanTime.toFixed(2) }}秒</span>
          </div>
        </el-card>
      </el-aside>

      <!-- 右侧文件列表 -->
      <el-main class="main-content">
        <el-card class="file-list-card">
          <template #header>
            <div class="card-header">
              <div class="header-left">
                <el-icon><Files /></el-icon>
                <span>文件列表</span>
                <el-tag v-if="filteredFiles.length" type="info" size="small">
                  {{ filteredFiles.length }} 个文件
                </el-tag>
              </div>
              <div class="header-right" v-if="filteredFiles.length > 0">
                <el-button
                  v-if="selectedFiles.length === 0"
                  type="primary"
                  plain
                  size="small"
                  @click="selectAllFiltered"
                >
                  <el-icon><Select /></el-icon>
                  全选 ({{ filteredFiles.length }})
                </el-button>
                <el-button
                  v-else
                  type="info"
                  plain
                  size="small"
                  @click="clearSelection"
                >
                  <el-icon><Close /></el-icon>
                  清除选择
                </el-button>
                <el-dropdown
                  split-button
                  type="success"
                  size="small"
                  @click="exportFiles"
                  @command="handleExportCommand"
                  :disabled="selectedFiles.length === 0"
                  style="margin-left: 8px;"
                >
                  <el-icon><Download /></el-icon>
                  导出选中 ({{ selectedFiles.length }})
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="folder">
                        <el-icon><Folder /></el-icon>
                        导出到文件夹
                      </el-dropdown-item>
                      <el-dropdown-item command="zip">
                        <el-icon><Files /></el-icon>
                        导出为压缩包
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
          </template>

          <!-- 空状态 -->
          <el-empty
            v-if="!scanResult"
            description="请选择路径并开始扫描"
          >
            <template #image>
              <el-icon :size="80" color="#909399"><FolderOpened /></el-icon>
            </template>
          </el-empty>

          <!-- 文件表格 -->
          <el-table
            v-else
            ref="tableRef"
            :data="pagedFiles"
            :row-key="(row: any) => row.path"
            style="width: 100%"
            class="file-table"
            @selection-change="handleSelectionChange"
            v-loading="scanning"
            border
          >
            <el-table-column type="selection" width="50" :reserve-selection="true" />

            <el-table-column label="状态" width="80">
              <template #default="scope">
                <el-tag v-if="scope.row.isValid === true" type="success" size="small">有效</el-tag>
                <el-tag v-else type="danger" size="small">无效</el-tag>
              </template>
            </el-table-column>

            <el-table-column label="类型" width="80">
              <template #default="scope">
                <el-tag :type="getTypeColor(scope.row.fileType)" size="small">
                  {{ (scope.row.fileType || '').toUpperCase() }}
                </el-tag>
              </template>
            </el-table-column>

            <el-table-column prop="name" label="文件名" width="300" />

            <el-table-column label="大小" width="100">
              <template #default="scope">
                {{ formatFileSize(scope.row.size) }}
              </template>
            </el-table-column>

            <el-table-column label="修改时间" width="160">
              <template #default="scope">
                {{ formatDate(scope.row.modTime) }}
              </template>
            </el-table-column>

            <el-table-column label="路径" min-width="400">
              <template #default="scope">
                <span class="path-link" @click.stop="openFolder(scope.row.path)">
                  {{ scope.row.path }}
                </span>
              </template>
            </el-table-column>
          </el-table>

          <!-- 分页 -->
          <div class="pagination-container" v-if="filteredFiles.length > 0">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[50, 100, 200, 500]"
              :total="filteredFiles.length"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handlePageChange"
              @current-change="handlePageChange"
            />
          </div>
        </el-card>
      </el-main>
    </el-container>

    <!-- 导出对话框 -->
    <el-dialog v-model="exportDialogVisible" :title="exportAsZip ? '导出为压缩包' : '导出文件'" width="500px">
      <el-form label-width="100px">
        <el-form-item label="导出目录">
          <div class="path-input">
            <el-input v-model="exportPath" readonly placeholder="选择导出目录" />
            <el-button @click="selectExportDirectory" type="primary">
              <el-icon><FolderOpened /></el-icon>
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="导出选项" v-if="!exportAsZip">
          <el-checkbox v-model="keepStructure">保持目录结构</el-checkbox>
        </el-form-item>
        <el-form-item label="" v-if="!exportAsZip">
          <el-checkbox v-model="overwriteExisting">覆盖已存在的文件</el-checkbox>
        </el-form-item>
        <el-form-item label="压缩选项" v-if="exportAsZip">
          <el-checkbox v-model="keepStructure">保持目录结构</el-checkbox>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="exportDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmExport" :loading="exporting">
          {{ exporting ? '导出中...' : (exportAsZip ? '导出压缩包' : '开始导出') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  GetDrives,
  SelectDirectory,
  SelectExportDirectory,
  ScanFiles,
  ExportFiles,
  ExportAsZip,
  FilterFiles,
  OpenFolder
} from '../wailsjs/go/main/App'
import { main, scanner } from '../wailsjs/go/models'
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime'

// 扫描进度类型
interface ScanProgressData {
  currentPath: string
  scannedDirs: number
  foundFiles: number
  currentFile: string
  isScanning: boolean
}

// 响应式状态
const drives = ref<main.DriveInfo[]>([])
const scanMode = ref('drive')
const selectedDrive = ref('')
const customPath = ref('')
const selectedTypes = ref(['pdf', 'word', 'excel', 'ppt'])
const validateFiles = ref(true)
const scanning = ref(false)
const scanResult = ref<scanner.ScanResult | null>(null)
const allFiles = ref<any[]>([])
const filteredFiles = ref<any[]>([])
const selectedFiles = ref<any[]>([])
const tableRef = ref<any>(null)

// 扫描进度
const scanProgress = ref<ScanProgressData>({
  currentPath: '',
  scannedDirs: 0,
  foundFiles: 0,
  currentFile: '',
  isScanning: false
})

// 过滤器状态
const filterText = ref('')
const validityFilter = ref('all')

// 分页状态
const currentPage = ref(1)
const pageSize = ref(100)

// 导出状态
const exportDialogVisible = ref(false)
const exportPath = ref('')
const keepStructure = ref(false)
const overwriteExisting = ref(false)
const exporting = ref(false)
const exportAsZip = ref(false)

// 初始化
onMounted(async () => {
  try {
    drives.value = await GetDrives()
    if (drives.value.length > 0) {
      selectedDrive.value = drives.value[0].path
    }
  } catch (error) {
    console.error('获取驱动器列表失败:', error)
  }

  // 监听扫描进度事件
  EventsOn('scan-progress', (progress: ScanProgressData) => {
    scanProgress.value = progress
  })
})

// 清理事件监听
onUnmounted(() => {
  EventsOff('scan-progress')
})

// 截断路径显示
const truncatePath = (path: string): string => {
  if (!path) return ''
  const maxLen = 40
  if (path.length <= maxLen) return path

  // 保留开头和结尾
  const start = path.substring(0, 15)
  const end = path.substring(path.length - 22)
  return `${start}...${end}`
}

// 选择目录
const selectDirectory = async () => {
  try {
    const path = await SelectDirectory()
    if (path) {
      customPath.value = path
    }
  } catch (error) {
    console.error('选择目录失败:', error)
  }
}

// 选择导出目录
const selectExportDirectory = async () => {
  try {
    const path = await SelectExportDirectory()
    if (path) {
      exportPath.value = path
    }
  } catch (error) {
    console.error('选择导出目录失败:', error)
  }
}

// 开始扫描
const startScan = async () => {
  const rootPath = scanMode.value === 'drive' ? selectedDrive.value : customPath.value

  if (!rootPath) {
    ElMessage.warning('请选择扫描路径')
    return
  }

  if (selectedTypes.value.length === 0) {
    ElMessage.warning('请至少选择一种文件类型')
    return
  }

  // 重置进度
  scanProgress.value = {
    currentPath: rootPath,
    scannedDirs: 0,
    foundFiles: 0,
    currentFile: '',
    isScanning: true
  }

  scanning.value = true
  try {
    const scanOptions = new scanner.ScanOptions({
      rootPath,
      includeTypes: selectedTypes.value,
      excludePaths: [],
      validateFiles: validateFiles.value
    })
    const result = await ScanFiles(scanOptions)

    scanResult.value = result
    allFiles.value = result.files || []
    filteredFiles.value = [...allFiles.value]
    currentPage.value = 1

    ElMessage.success(`扫描完成，共找到 ${result.totalCount} 个文件`)
  } catch (error: any) {
    ElMessage.error('扫描失败: ' + error.message)
  } finally {
    scanning.value = false
    scanProgress.value.isScanning = false
  }
}

// 应用过滤器 - 在前端进行过滤，避免频繁调用后端
const applyFilter = () => {
  if (!allFiles.value.length) {
    filteredFiles.value = []
    return
  }

  const searchText = filterText.value.toLowerCase()
  const validOnly = validityFilter.value === 'valid'
  const invalidOnly = validityFilter.value === 'invalid'

  filteredFiles.value = allFiles.value.filter(file => {
    // 按有效性过滤
    if (validOnly && !file.isValid) return false
    if (invalidOnly && file.isValid) return false

    // 按文件名搜索
    if (searchText && !file.name.toLowerCase().includes(searchText)) {
      return false
    }

    return true
  })

  // 重置到第一页
  currentPage.value = 1
}

// 计算当前页的数据
const pagedFiles = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredFiles.value.slice(start, end)
})

// 计算是否已全选所有过滤后的文件
const allFilteredSelected = computed(() => {
  return filteredFiles.value.length > 0 &&
         selectedFiles.value.length === filteredFiles.value.length
})

// 全选所有过滤后的文件（跨页）
const selectAllFiltered = () => {
  if (!tableRef.value) return

  // 先清除所有选择
  tableRef.value.clearSelection()

  // 选中所有过滤后的文件
  filteredFiles.value.forEach((file: any) => {
    tableRef.value.toggleRowSelection(file, true)
  })

  ElMessage.success(`已选中 ${filteredFiles.value.length} 个文件`)
}

// 清除所有选择
const clearSelection = () => {
  if (tableRef.value) {
    tableRef.value.clearSelection()
  }
  selectedFiles.value = []
}

// 处理分页变化
const handlePageChange = () => {
  // 使用 reserve-selection 后，分页变化时会自动保持选择状态
}

// 处理选择变化
const handleSelectionChange = (selection: any[]) => {
  selectedFiles.value = selection
}

// 导出文件（默认导出到文件夹）
const exportFiles = () => {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning('请先选择要导出的文件')
    return
  }
  exportAsZip.value = false
  exportDialogVisible.value = true
}

// 处理导出命令
const handleExportCommand = (command: string) => {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning('请先选择要导出的文件')
    return
  }
  exportAsZip.value = command === 'zip'
  exportDialogVisible.value = true
}

// 打开文件所在文件夹
const openFolder = async (filePath: string) => {
  try {
    await OpenFolder(filePath)
  } catch (error) {
    console.error('打开文件夹失败:', error)
    ElMessage.error('打开文件夹失败')
  }
}

// 确认导出
const confirmExport = async () => {
  if (!exportPath.value) {
    ElMessage.warning('请选择导出目录')
    return
  }

  exporting.value = true
  try {
    const exportOptions = {
      destPath: exportPath.value,
      files: selectedFiles.value,
      keepStructure: keepStructure.value,
      overwrite: overwriteExisting.value
    }

    let result
    if (exportAsZip.value) {
      result = await ExportAsZip(exportOptions as any)
    } else {
      result = await ExportFiles(exportOptions as any)
    }

    exportDialogVisible.value = false

    if (result.failed > 0) {
      ElMessageBox.alert(
        `成功导出 ${result.success} 个文件，失败 ${result.failed} 个`,
        '导出完成',
        { type: 'warning' }
      )
    } else {
      const msg = exportAsZip.value ? `成功导出 ${result.success} 个文件到压缩包` : `成功导出 ${result.success} 个文件`
      ElMessage.success(msg)
    }
  } catch (error: any) {
    ElMessage.error('导出失败: ' + (error.message || error))
  } finally {
    exporting.value = false
  }
}

// 工具函数 - 安全获取文件属性（处理可能的大小写差异）
const getFileValid = (row: any): boolean => {
  // 尝试多种可能的属性名
  if (typeof row.isValid === 'boolean') return row.isValid
  if (typeof row.IsValid === 'boolean') return row.IsValid
  return true // 默认有效
}

const getInvalidReason = (row: any): string => {
  return row.invalidReason || row.InvalidReason || '未知原因'
}

const getFileType = (row: any): string => {
  return row.fileType || row.FileType || row.extension?.replace('.', '') || ''
}

const getFileSize = (row: any): number => {
  return row.size || row.Size || 0
}

const getTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    pdf: 'danger',
    word: 'primary',
    excel: 'success',
    ppt: 'warning'
  }
  return colors[type?.toLowerCase()] || 'info'
}

const formatFileSize = (bytes: number) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}
</script>

<style scoped>
.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  padding: 0 24px;
  height: 60px;
}

.title {
  margin: 0;
  font-size: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.main-container {
  flex: 1;
  overflow: hidden;
}

.sidebar {
  background: #fff;
  padding: 16px;
  overflow-y: auto;
  border-right: 1px solid #e4e7ed;
}

.control-card {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-right {
  margin-left: auto;
}

.form-section {
  margin-bottom: 16px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  color: #606266;
}

.scan-mode-group {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.path-selector {
  margin-top: 8px;
}

.full-width {
  width: 100%;
}

.path-input {
  display: flex;
  gap: 8px;
}

.path-input .el-input {
  flex: 1;
}

.scan-btn {
  margin-top: 8px;
}

.scan-progress {
  margin-top: 16px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
}

.progress-stats {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 12px;
  color: #606266;
}

.progress-bar-container {
  width: 100%;
  height: 6px;
  background: #e4e7ed;
  border-radius: 3px;
  overflow: hidden;
}

.progress-bar-animated {
  width: 30%;
  height: 100%;
  background: linear-gradient(90deg, #409eff, #66b1ff, #409eff);
  background-size: 200% 100%;
  border-radius: 3px;
  animation: progress-flow 1.5s ease-in-out infinite, progress-move 2s ease-in-out infinite;
}

@keyframes progress-flow {
  0% {
    background-position: 100% 0;
  }
  100% {
    background-position: -100% 0;
  }
}

@keyframes progress-move {
  0% {
    transform: translateX(-100%);
  }
  50% {
    transform: translateX(233%);
  }
  100% {
    transform: translateX(-100%);
  }
}

.current-path {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.current-path .el-icon {
  flex-shrink: 0;
  color: #409eff;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  margin: 12px 0;
}

.stat-item {
  text-align: center;
  padding: 8px 4px;
  background: #f5f7fa;
  border-radius: 6px;
}

.stat-value {
  display: block;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.stat-value.valid {
  color: #67c23a;
}

.stat-value.invalid {
  color: #f56c6c;
}

.stat-value.selected {
  color: #409eff;
}

.stat-label {
  font-size: 11px;
  color: #909399;
}

.scan-time {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #909399;
  padding-top: 8px;
  border-top: 1px solid #ebeef5;
}

.path-link {
  color: #409eff;
  cursor: pointer;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: block;
}

.path-link:hover {
  text-decoration: underline;
}

.main-content {
  padding: 16px;
  overflow: hidden;
}

.file-list-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.file-list-card :deep(.el-card__header) {
  flex-shrink: 0;
}

.file-list-card :deep(.el-card__body) {
  padding: 12px;
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.file-table {
  flex: 1;
  min-height: 0;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: center;
  flex-shrink: 0;
}

:deep(.el-checkbox-group) {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

:deep(.el-checkbox) {
  margin-right: 0;
}

:deep(.el-radio-group) {
  width: 100%;
}

:deep(.el-radio-button) {
  flex: 1;
}

:deep(.el-radio-button__inner) {
  width: 100%;
}
</style>
