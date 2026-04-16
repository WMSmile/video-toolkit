<template>
  <el-card>
    <template #header>
      <div class="flex justify-between items-center">
        <h2 class="text-xl font-bold">视频转MP3</h2>
      </div>
    </template>

    <el-col>
      <el-form label-width="80px">
        <el-row>
          <el-col :span="24">
            <el-form-item label="选择文件">
              <el-col :span="24" class="flex">
                <el-input v-model="filePathsText" placeholder="请选择视频文件" readonly>
                </el-input>
                <el-button type="primary" @click="chooseSingleFile">
                  单选
                </el-button>
                <el-button type="success" @click="chooseMultipleFiles" class="ml-2">
                  多选
                </el-button>
              </el-col>
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item label="歌名">
              <el-input v-model="title" placeholder="请输入歌名（可选）" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="歌手">
              <el-input v-model="artist" placeholder="请输入歌手（可选）" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item>
              <el-button type="primary" @click="start" :loading="isLoading">
                {{ isLoading ? '转换中...' : '开始转换' }}
              </el-button>
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item v-if="isLoading">
              <el-progress :percentage="progress" :status="progressStatus" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item v-if="status">
              <el-alert :title="status" :type="statusType" show-icon :closable="false" />
            </el-form-item>
          </el-col>
          <el-col :span="24" v-if="conversionResults.length > 0">
            <el-form-item label="转换结果">
              <el-table :data="conversionResults" style="width: 100%">
                <el-table-column prop="file" label="文件" width="400" />
                <el-table-column prop="status" label="状态" />
                <el-table-column prop="output" label="输出文件" />
              </el-table>
            </el-form-item>
          </el-col>

        </el-row>

      </el-form>
    </el-col>

  </el-card>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

// 响应式数据
const title = ref('')
const artist = ref('')
const status = ref('')
const filePaths = ref([])
const filePathsText = ref('')
const isLoading = ref(false)
const progress = ref(0)
const progressStatus = ref('')
const convertId = ref('')
const conversionResults = ref([])

// 计算属性
const statusType = computed(() => {
  if (status.value.includes('完成')) return 'success'
  if (status.value.includes('转换中')) return 'info'
  if (status.value.includes('错误')) return 'error'
  return 'info'
})

// 导入服务和事件
import { ConvertService, FileService, ID3Service } from "../../bindings/video-toolkit/services";
import { Events } from "@wailsio/runtime";

// 生成唯一的转换ID
function generateConvertId() {
  return 'conv_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
}

// 选择单个文件
async function chooseSingleFile() {
  const path = await FileService.SelectFile("请选择待处理视频", "视频文件", "*.mp4;*.mov;*.mkv;*.avi")
  if (path === "") {
    console.log("用户取消了选择")
    return
  }
  console.log("选中的视频路径:", path)
  filePaths.value = [path]
  filePathsText.value = path
  status.value = ''
  progress.value = 0
  conversionResults.value = []
}

// 选择多个文件
async function chooseMultipleFiles() {
  const paths = await FileService.SelectMultipleFiles("请选择待处理视频")
  if (paths && paths.length > 0) {
    console.log("选中的视频路径:", paths)
    filePaths.value = paths
    filePathsText.value = paths.join(', ')
    status.value = ''
    progress.value = 0
    conversionResults.value = []
  }
}

// 进度回调函数
function progressCallback(value) {
  progress.value = value
  console.log('进度更新:', value)
}

// 开始转换
async function start() {
  if (filePaths.value.length === 0) {
    status.value = '错误：请选择视频文件'
    return
  }

  // 重置状态
  isLoading.value = true
  status.value = '转换中...'
  progress.value = 0
  progressStatus.value = ''
  conversionResults.value = []

  try {
    // 遍历所有选择的文件
    let successCount = 0
    let totalCount = filePaths.value.length

    for (let i = 0; i < totalCount; i++) {
      const filePath = filePaths.value[i]
      // 生成转换ID
      convertId.value = generateConvertId()

      // 更新进度
      progress.value = Math.round((i / totalCount) * 100)
      status.value = `转换中... (${i + 1}/${totalCount})`

      try {
        // 调用转换服务（带进度回调）
        console.log('开始转换:', filePath, '转换ID:', convertId.value)
        const out = await ConvertService.ToMP3Progress(filePath, convertId.value, progressCallback)
        console.log('转换完成:', out)

        // 写入ID3标签
        if (title.value || artist.value) {
          await ID3Service.Write(out, { Title: title.value, Artist: artist.value, Album: '' })
        }

        // 添加到结果列表
        conversionResults.value.push({
          file: filePath,
          status: '成功',
          output: out
        })

        successCount++
      } catch (error) {
        console.error('转换失败:', error)
        // 添加到结果列表
        conversionResults.value.push({
          file: filePath,
          status: `失败: ${error.message || '转换失败'}`,
          output: ''
        })
      }
    }

    progress.value = 100
    progressStatus.value = 'success'

    // 显示最终状态
    if (successCount === totalCount) {
      status.value = `✅ 全部完成：成功转换 ${successCount} 个文件`
    } else if (successCount > 0) {
      status.value = `⚠️ 部分完成：成功转换 ${successCount} 个文件，失败 ${totalCount - successCount} 个文件`
    } else {
      status.value = `❌ 全部失败：所有 ${totalCount} 个文件转换失败`
    }
  } catch (error) {
    console.error('批量转换失败:', error)
    status.value = `错误：${error.message || '批量转换失败'}`
    progressStatus.value = 'exception'
  } finally {
    isLoading.value = false
  }
}

// 监听进度事件
function handleProgressEvent(data) {
  try {
    const parsedData = JSON.parse(data)
    console.log('收到进度事件:', parsedData)

    // 只处理当前转换的进度
    if (parsedData.convertId === convertId.value) {
      progress.value = parsedData.progress
    }
  } catch (error) {
    console.error('解析进度事件失败:', error)
  }
}

// 生命周期钩子
onMounted(() => {
  // 监听进度事件
  Events.On("convert_progress", handleProgressEvent)
  console.log('已注册进度事件监听器')
})

onUnmounted(() => {
  // 移除进度事件监听器
  Events.Off("convert_progress", handleProgressEvent)
  console.log('已移除进度事件监听器')
})
</script>

<style scoped>
.flex {
  display: flex;
}

.justify-between {
  justify-content: space-between;
}

.items-center {
  align-items: center;
}

.text-xl {
  font-size: 1.25rem;
}

.font-bold {
  font-weight: bold;
}
</style>