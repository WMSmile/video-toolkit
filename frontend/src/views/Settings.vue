<template>
  <el-card>
    <h2>⚙ 设置</h2>
    <el-tabs v-model="activeTab" class="mt-4">
      <!-- FFmpeg 配置 -->
      <el-tab-pane label="FFmpeg 配置" name="ffmpeg">
        <el-form label-width="100px">
          <el-form-item label="FFmpeg 路径">
            <el-input v-model="ffmpegPath" placeholder="请输入 FFmpeg 可执行文件路径,默认ffmpeg" />
          </el-form-item>
          <el-form-item label="FFprobe 路径">
            <el-input v-model="ffprobePath" placeholder="请输入 FFprobe 可执行文件路径,默认ffprobe" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="saveFFmpegConfig">保存配置</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- 主题配置 -->
      <el-tab-pane label="主题配置" name="theme">
        <el-form label-width="100px">
          <el-form-item label="主题模式">
            <el-select v-model="theme" placeholder="选择主题" @change="changeTheme">
              <el-option label="浅色" value="light" />
              <el-option label="深色" value="dark" />
              <el-option label="蓝色" value="blue" />
            </el-select>
          </el-form-item>

          <el-form-item label="主色调">
            <el-color-picker v-model="primaryColor" @change="changePrimaryColor" />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- sherpa-onnx模型配置 -->
      <el-tab-pane label="sherpa-onnx模型配置" name="sherpa">
        <el-form label-width="100px">
          <el-form-item label="ModelPath">
            <el-input v-model="modelPath" placeholder="请输入模型路径" />
          </el-form-item>
          <el-form-item label="TokensPath">
            <el-input v-model="tokensPath" placeholder="请输入词表路径" />
          </el-form-item>
          <el-form-item label="VadModelPath">
            <el-input v-model="vadModelPath" placeholder="请输入VAD模型路径" />
          </el-form-item>
          <el-form-item label="PunctuationModelPath">
            <el-input v-model="punctuationModelPath" placeholder="请输入标点模型路径" />
          </el-form-item>
          <el-form-item label="HotwordsPath">
            <el-input v-model="hotwordsPath" placeholder="请输入热词路径" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="saveSherpaConfig">保存配置</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- 关于 -->
      <el-tab-pane label="关于" name="about">
        <div class="about-content">
          <h3>视频工具箱</h3>
          <p>版本: 1.0.0</p>
          <p>这是一个使用 Vue 3 + Element Plus + Vue Router 的桌面端应用，用于处理视频文件。</p>
          <p>支持视频格式转换、批量处理、视频标签编辑等功能。</p>
          <p>项目使用了 Silero VAD 模型和 Punct-CT Transformer 模型。</p>
          <p>作者：武猛</p>
          <p class="mt-4">© 2026 视频工具箱</p>
        </div>
      </el-tab-pane>
    </el-tabs>
  </el-card>
</template>

<script setup name="Settings">
import { ref, onMounted } from 'vue'

// 导入 ElMessage
import { ElMessage } from 'element-plus'
import { lo } from 'element-plus/es/locale/index.mjs';

// 导入服务和事件
import { ConfigService } from "../../bindings/video-toolkit/services";
import { Events } from "@wailsio/runtime";

// 响应式数据
const activeTab = ref('ffmpeg')
const ffmpegPath = ref('')
const ffprobePath = ref('')
const theme = ref('light')
const primaryColor = ref(
  getComputedStyle(document.documentElement)
    .getPropertyValue('--el-color-primary')
    .trim()
)
const modelPath = ref('')
const tokensPath = ref('')
const vadModelPath = ref('')
const punctuationModelPath = ref('')
const hotwordsPath = ref('')

// 生命周期钩子
onMounted(() => {
  // 加载保存的配置
  loadConfig()

  
})

async function loadConfig() {
  // 从本地存储加载配置
  const config = await ConfigService.GetConfig()
  console.log(config)
  ffmpegPath.value = config.FFmpegPath
  ffprobePath.value = config.FFprobePath
  theme.value = config.Theme
  primaryColor.value = config.PrimaryColor
  modelPath.value = config.ModelPath
  tokensPath.value = config.TokensPath
  vadModelPath.value = config.VadModelPath
  punctuationModelPath.value = config.PunctuationModelPath
  hotwordsPath.value = config.HotwordsPath

  changeTheme(theme.value)
  changePrimaryColor(primaryColor.value)
}

// 方法
async function saveFFmpegConfig() {
  let config = await ConfigService.GetConfig()
  config.FFmpegPath = ffmpegPath.value
  config.FFprobePath = ffprobePath.value
  // 保存到本地存储
  ConfigService.SetConfig(config)
  // 提示保存成功
  ElMessage.success('FFmpeg 配置保存成功')
}

async function saveSherpaConfig() {
  let config = await ConfigService.GetConfig()
  config.ModelPath = modelPath.value
  config.TokensPath = tokensPath.value
  config.VadModelPath = vadModelPath.value
  config.PunctuationModelPath = punctuationModelPath.value
  config.HotwordsPath = hotwordsPath.value
  // 保存到本地存储
  ConfigService.SetConfig(config)
  // 提示保存成功
  ElMessage.success('sherpa-onnx模型配置保存成功')
}



async function changeTheme(val) {
  // 移除所有现有的主题类
  document.documentElement.classList.remove('light', 'dark', 'blue')
  // 添加新的主题类
  if (val === 'dark' || val === 'light' || val === 'blue') {
    document.documentElement.classList.add(val)
    // 保存主题设置
    let config = await ConfigService.GetConfig()
    config.Theme = val
    // 保存到本地存储
    ConfigService.SetConfig(config)

  } else {
    console.error('未知的主题模式')
  }
}

async function changePrimaryColor(color) {
  document.documentElement.style.setProperty('--el-color-primary', color)
  // 保存主色调设置
  let config = await ConfigService.GetConfig()
  config.PrimaryColor = color
  // 保存到本地存储
  ConfigService.SetConfig(config)
}


</script>

<style scoped>
.about-content {
  padding: 20px;
  line-height: 1.6;
}

.mt-4 {
  margin-top: 16px;
}
</style>