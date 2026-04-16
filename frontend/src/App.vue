<template>
  <el-container class="app-container">
    <!-- 左侧菜单 -->
    <el-aside width="200px" class="sidebar">
      <el-text class="logo">🪟 Video Toolkit</el-text>
      <el-menu :default-active="$route.path" router>
        <el-menu-item index="/"><el-icon>
            <HomeFilled />
          </el-icon>首页</el-menu-item>
        <el-menu-item index="/settings"><el-icon>
            <SettingsFilled />
          </el-icon>设置</el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 右侧内容 -->
    <el-main class="content">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup name="App">
// App 组件
</script>

<style scoped lang="scss">
.app-container {
  height: 100vh;
  width: 100vw;
  background-color: var(--el-bg-color-page);
}

.sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding-top: 10px;
}

.logo {
  font-size: 18px;
  font-weight: bold;
  padding: 15px;
  text-align: center;
}

.content {
  background-color: var(--el-bg-color-page);
  padding: 20px;
}
</style>