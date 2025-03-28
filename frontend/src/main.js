import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/antd.css'
import './styles/index.css'
import { setupErrorHandler } from './utils/error-handler'
import { setupPermission } from './directives/permission'
import { setupGlobalComponents } from './components'
import { useSettingsStore } from '@/store/modules/settings'

const app = createApp(App) // 创建Vue应用实例
const pinia = createPinia() // 创建状态管理实例

app.use(pinia) // 使用Pinia状态管理

// 初始化设置
const settingsStore = useSettingsStore()
settingsStore.initSettings()

// 注册全局错误处理
setupErrorHandler(app)

// 注册全局组件
setupGlobalComponents(app)

// 注册权限指令
setupPermission(app)

app.use(router) // 使用路由
app.use(Antd) // 使用Ant Design Vue组件库

app.mount('#app') // 将应用挂载到 #app 节点