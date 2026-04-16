<template>
  <el-card>
    <template #header>
      <div class="flex justify-between items-center">
        <h2 class="text-xl font-bold">MP3标签编辑</h2>
      </div>
    </template>

    <el-col>
      <el-form label-width="80px">
        <el-row>
          <el-col :span="24">
            <el-form-item label="选择文件">
              <el-col :span="24" class="flex">
                <el-input v-model="filePath" placeholder="请选择MP3文件" readonly>
                </el-input>
                <el-button type="primary" @click="chooseFile">
                  选择文件
                </el-button>
              </el-col>
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item label="标题">
              <el-input v-model="tag.Title" placeholder="请输入标题" />
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item label="歌手">
              <el-input v-model="tag.Artist" placeholder="请输入歌手" />
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item label="专辑">
              <el-input v-model="tag.Album" placeholder="请输入专辑" />
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item>
              <el-button type="primary" @click="load" :loading="isLoading">
                {{ isLoading ? '读取中...' : '读取标签' }}
              </el-button>
              <el-button type="success" @click="save" :loading="isLoading" class="ml-2">
                {{ isLoading ? '保存中...' : '保存标签' }}
              </el-button>
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item v-if="status">
              <el-alert :title="status" :type="statusType" show-icon :closable="false" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-col>
  </el-card>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ID3Service, FileService } from "../../bindings/video-toolkit/services";

// 响应式数据
const tag = ref({ Title: '', Artist: '', Album: '' })
const filePath = ref('')
const currentFile = ref('')
const isLoading = ref(false)
const status = ref('')

// 计算属性
const statusType = computed(() => {
  if (status.value.includes('成功')) return 'success'
  if (status.value.includes('错误')) return 'error'
  return 'info'
})

// 选择文件
async function chooseFile() {
  const path = await FileService.SelectFile("请选择MP3文件", "音频文件", "*.mp3")
  if (path === "") {
    console.log("用户取消了选择")
    return
  }
  console.log("选中的MP3路径:", path)
  filePath.value = path
  currentFile.value = path
  status.value = ''
}

// 读取标签
async function load() {
  if (!currentFile.value) {
    status.value = '错误：请选择MP3文件'
    return
  }

  isLoading.value = true
  status.value = '读取标签中...'

  try {
    const res = await ID3Service.Read(currentFile.value)
    console.log(res)
    tag.value = res
    console.log(tag.value)
    status.value = '✅ 标签读取成功'
  } catch (error) {
    console.error('读取标签失败:', error)
    status.value = `错误：${error.message || '读取标签失败'}`
  } finally {
    isLoading.value = false
  }
}

// 保存标签
async function save() {
  if (!currentFile.value) {
    status.value = '错误：请选择MP3文件'
    return
  }

  isLoading.value = true
  status.value = '保存标签中...'

  try {
    await ID3Service.Write(currentFile.value, tag.value)
    status.value = '✅ 标签保存成功'
  } catch (error) {
    console.error('保存标签失败:', error)
    status.value = `错误：${error.message || '保存标签失败'}`
  } finally {
    isLoading.value = false
  }
}
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

.ml-2 {
  margin-left: 8px;
}
</style>