import { createApp } from 'vue'
import App from './App.vue'

// 引入 Element Plus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

// 引入路由
import router from './router'

// 引入主题样式
import './styles/theme.scss'


// 创建 Vue 实例并挂载应用
createApp(App)
.use(ElementPlus)
.use(router)
.mount('#app')
