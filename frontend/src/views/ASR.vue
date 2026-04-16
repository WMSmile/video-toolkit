<template>
  <el-card>
    <template #header>
      <div class="flex justify-between items-center">
        <h2 class="text-xl font-bold">批量语音识别</h2>
      </div>
    </template>

    <el-col>
      <el-form label-width="80px">
        <el-row>
          <el-col :span="24">
            <el-form-item label="选择文件">
              <el-col :span="24" class="flex">
                <el-input v-model="filePathsText" placeholder="请选择音频文件" readonly>
                </el-input>
                <el-button type="primary" @click="chooseFiles">
                  选择文件
                </el-button>
              </el-col>
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item>
              <el-button type="primary" @click="startRecognize" :loading="isLoading">
                {{ isLoading ? '识别中...' : '开始识别' }}
              </el-button>
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item v-if="status">
              <el-alert :title="status" :type="statusType" show-icon :closable="false" />
            </el-form-item>
          </el-col>

          <el-col :span="24" v-if="results.length > 0">
            <el-form-item label="识别结果">
              <el-col :span="24" class="flex">
                <el-input
                  v-model="outputText"
                  type="textarea"
                  :rows="10"
                  placeholder="识别结果将显示在这里"
                  readonly
                />
                <el-button
                  type="primary"
                  @click="copyOutput"
                  class="ml-2"
                  :disabled="!outputText"
                >
                  复制
                </el-button>
              </el-col>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-col>
  </el-card>
</template>

<script setup>
import { ref, computed } from 'vue'
import { FileService, ASRService } from "../../bindings/video-toolkit/services";

// 响应式数据
const filePaths = ref([])
const filePathsText = ref('')
const isLoading = ref(false)
const status = ref('')
const results = ref([])
const outputText = ref('')

// 计算属性
const statusType = computed(() => {
  if (status.value.includes('成功')) return 'success'
  if (status.value.includes('错误')) return 'error'
  return 'info'
})

// 选择多个文件
async function chooseFiles() {
  const paths = await FileService.SelectMultipleFiles("请选择音频文件", "音频文件", "*.mp3;*.wav;*.flac;*.aac;*.m4a;*.ogg;*.mp4;*.flv;*.mkv;*.avi;*.mov;*.wmv")
  if (paths && paths.length > 0) {
    filePaths.value = paths
    filePathsText.value = paths.join(', ')
    status.value = ''
    results.value = []
  }
}

// 开始识别
async function startRecognize() {
  if (filePaths.value.length === 0) {
    status.value = '错误：请选择音频文件'
    return
  }

  isLoading.value = true
  status.value = '识别中...'

  try {
    const res = await ASRService.BatchRecognize(filePaths.value)
    
    // 转换结果格式
    results.value = Object.entries(res).map(([file, text]) => ({
      file,
      text
    }))
    
    // 生成输出文本
    outputText.value = results.value.map(item => {
      const fileName = item.file.split('/').pop()
      return `${fileName}\n${item.text}\n`
    }).join('')
    
    status.value = '✅ 识别完成'
  } catch (error) {
    console.error('识别失败:', error)
    status.value = `错误：${error.message || '识别失败'}`
  } finally {
    isLoading.value = false
  }
}

// 复制输出文本
function copyOutput() {
  if (outputText.value) {
    navigator.clipboard.writeText(outputText.value)
      .then(() => {
        ElMessage.success('复制成功')
      })
      .catch(err => {
        console.error('复制失败:', err)
        ElMessage.error('复制失败')
      })
  }
}

// 导入 ElMessage
import { ElMessage } from 'element-plus'
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