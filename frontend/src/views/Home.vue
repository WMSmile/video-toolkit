<template>
  <el-card>
    <h2>🏠 首页</h2>
    <p>欢迎来到 Windows 桌面风格 Vue 应用</p>

    <!-- 九宫格card -->
    <el-card-grid :gutter="20">
      <el-card-grid-item :span="6">
        <router-link to="/convert">
          <el-card>
            <h2>视频转MP3</h2>
            <p>将视频转换为MP3格式</p>
          </el-card>
        </router-link>
      </el-card-grid-item>
      <el-card-grid-item :span="6">
        <router-link to="/tagger">
          <el-card>
            <h2>视频标签</h2>
            <p>为视频添加标签</p>
          </el-card>
        </router-link>
      </el-card-grid-item>

      <el-card-grid-item :span="6">
        <router-link to="/asr">
          <el-card>
            <h2>批量语音识别</h2>
            <p>批量识别音频文件内容</p>
          </el-card>
        </router-link>
      </el-card-grid-item>
      </el-card-grid>
    </el-card>

</template>

<script setup name="Home">
// Home 组件
</script>