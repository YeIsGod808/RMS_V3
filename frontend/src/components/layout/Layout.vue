<template>
  <a-config-provider :theme="settingsStore.themeConfig">
    <div class="layout">
      <a-layout class="layout">
        <a-layout-header class="header">
          <!-- Logo区域 -->
          <div class="logo" @click="router.push('/')">
            <database-outlined class="logo-icon" />
            <span class="logo-text">数据库教学资源管理平台</span>
          </div>

          <a-menu
            v-model:selectedKeys="selectedKeys"
            theme="dark"
            mode="horizontal"
            :style="{ lineHeight: '80px', flex: 1 }"
            class="main-menu"
          >
            <a-menu-item key="graph" class="menu-item" style="order: 1;">
              <router-link to="/graph">知识图谱</router-link>
            </a-menu-item>
            <a-menu-item key="route" class="menu-item" style="order: 2;">
              <router-link to="/route">学习规划</router-link>
            </a-menu-item>
            <a-menu-item key="outside" class="menu-item" style="order: 3;">
              <router-link to="/external">外部资源</router-link>
            </a-menu-item>
            <div v-permission="'admin'" style="order: 4;">
              <a-menu-item key="admin" class="menu-item">
                <router-link to="/admin">管理后台</router-link>
              </a-menu-item>
            </div>
          </a-menu>
          
          <div class="header-right">
            <a-space>
              <!-- 主题切换 -->
              <a-dropdown>
                <a-button type="text" class="header-btn">
                  <template #icon><bg-colors-outlined /></template>
                  主题
                </a-button>
                <template #overlay>
                  <a-menu class="theme-menu">
                    <a-menu-item key="default" @click="settingsStore.setTheme('default')">
                      <check-outlined v-if="settingsStore.theme === 'default'" />
                      默认主题
                    </a-menu-item>
                    <!-- <a-menu-item key="dark" @click="settingsStore.setTheme('dark')">
                      <check-outlined v-if="settingsStore.theme === 'dark'" />
                      暗色主题
                    </a-menu-item>
                    <a-menu-item key="compact" @click="settingsStore.setTheme('compact')">
                      <check-outlined v-if="settingsStore.theme === 'compact'" />
                      紧凑主题
                    </a-menu-item> -->
                  </a-menu>
                </template>
              </a-dropdown>

              <!-- 用户信息 -->
              <a-dropdown>
                <a-button type="text" class="header-btn user-btn">
                  <template #icon><user-outlined /></template>
                  {{ userStore.nickname || userStore.username }}
                  <down-outlined />
                </a-button>
                <template #overlay>
                  <a-menu class="user-menu">
                    <a-menu-item>
                      <user-outlined />
                      <span>{{ userStore.nickname || userStore.username }}</span>
                    </a-menu-item>
                    <a-menu-item>
                      <team-outlined />
                      <span>{{ getRoleText(userStore.role) }}</span>
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item @click="showChangePassword">
                      <key-outlined />
                      修改密码
                    </a-menu-item>
                    <a-menu-item @click="handleLogout">
                      <logout-outlined />
                      退出登录
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </a-space>
          </div>
        </a-layout-header>

        <a-layout-content class="content">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </a-layout-content>
        
        <!-- 页脚 -->
        <a-layout-footer class="footer">
          <div class="footer-content">
            <div class="copyright">
              © {{ new Date().getFullYear() }} 数据库教学资源管理平台
            </div>
          </div>
        </a-layout-footer>
      </a-layout>
    </div>

    <!-- 修改密码组件 -->
    <change-password v-model:visible="changePasswordVisible" @success="handlePasswordChanged" />
  </a-config-provider>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/modules/user'
import { useSettingsStore } from '@/store/modules/settings'
import { message } from 'ant-design-vue'
import {
  UserOutlined,
  LogoutOutlined,
  BgColorsOutlined,
  CheckOutlined,
  DownOutlined,
  TeamOutlined,
  KeyOutlined,
  DatabaseOutlined  // 添加数据库图标
} from '@ant-design/icons-vue'
import ChangePassword from '@/components/user/ChangePassword.vue'

const router = useRouter()
const userStore = useUserStore()
const settingsStore = useSettingsStore()
const selectedKeys = ref()
const changePasswordVisible = ref(false)

const handleLogout = async () => {
  try {
    await userStore.logout()
    // 使用nextTick确保状态更新后再跳转
    nextTick(() => {
      router.push('/login')
    })
  } catch (error) {
    console.error(error)
  }
}

// 获取角色文本
const getRoleText = (role) => {
  const roleMap = {
    'admin': '管理员',
    'teacher': '教师',
    'student': '学生'
  }
  return roleMap[role] || '用户'
}

// 显示修改密码弹窗
const showChangePassword = async () => {
  try {
    // 先验证token有效性
    const isValid = await userStore.verifyToken()
    if (isValid) {
      changePasswordVisible.value = true
    } else {
      message.error('登录已过期，请重新登录')
      await userStore.logout()
      router.push('/login')
    }
  } catch (error) {
    console.error('验证token失败:', error)
    message.error('验证身份信息失败，请重新登录')
    await userStore.logout()
    router.push('/login')
  }
}

// 密码修改成功的回调
const handlePasswordChanged = () => {
  message.success('密码已更新')
}
</script>

<style scoped>
/* 全局布局样式 */
.layout {
  min-height: 100vh;
  background: var(--background-color, #f5f7fa);
}

/* 顶部导航栏 */
.header {
  display: flex;
  align-items: center;
  padding: 0 24px;
  background: linear-gradient(90deg, #001529 0%, #003a70 100%, #004a8c 120%);
  height: 80px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  position: relative;
  z-index: 10;
  transition: all 0.3s ease;
}

/* Logo区域 */
.logo {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 24px;
  transition: all 0.3s ease-out;
  border-radius: 6px;
  margin-right: 20px;
  position: relative;
  overflow: hidden;
}

.logo::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0);
  transition: all 0.3s ease;
  z-index: -1;
}

.logo:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.logo:hover::after {
  background: rgba(255, 255, 255, 0.05);
}

.logo-icon {
  font-size: 30px;
  color: #1890ff;
  margin-right: 14px;
  filter: drop-shadow(0 0 5px rgba(24, 144, 255, 0.6));
  transition: all 0.4s ease;
}

.logo:hover .logo-icon {
  transform: rotate(15deg) scale(1.1);
  color: #40a9ff;
}

.logo-text {
  color: white;
  font-size: 20px;
  font-weight: 600;
  background: linear-gradient(120deg, #1890ff, #69c0ff, #91d5ff);
  background-size: 200% auto;
  -webkit-background-clip: text;
  /* 添加标准属性 */
  background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: 0 0 12px rgba(24, 144, 255, 0.4);
  transition: all 0.5s ease;
  letter-spacing: 0.5px;
}

.logo:hover .logo-text {
  background-position: right center;
}

/* 主菜单 */
.main-menu {
  border-bottom: none;
  margin-left: 10px;
}

.menu-item {
  position: relative;
  margin: 0 8px;
}

:deep(.ant-menu-horizontal) {
  border-bottom: none;
  background: transparent;
}

:deep(.ant-menu-item) {
  border-radius: 6px;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  margin: 0 5px;
  padding: 0 18px;
  position: relative;
  overflow: hidden;
}

:deep(.ant-menu-item::before) {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(120deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  transform: translateX(-100%);
  transition: all 0.6s ease;
}

:deep(.ant-menu-item:hover::before) {
  transform: translateX(100%);
}

:deep(.ant-menu-item:hover) {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

:deep(.ant-menu-item a) {
  color: rgba(255, 255, 255, 0.85);
  font-weight: 500;
  font-size: 16px;
  letter-spacing: 0.5px;
  transition: all 0.3s;
  position: relative;
  padding: 0 4px;
}

:deep(.ant-menu-item a::after) {
  content: '';
  position: absolute;
  bottom: -5px;
  left: 0;
  width: 100%;
  height: 2px;
  background: #1890ff;
  transform: scaleX(0);
  transform-origin: right;
  transition: transform 0.4s ease;
}

:deep(.ant-menu-item:hover a::after) {
  transform: scaleX(1);
  transform-origin: left;
}

:deep(.ant-menu-item-selected) {
  background: rgba(24, 144, 255, 0.25) !important;
}

:deep(.ant-menu-item-selected a) {
  color: white !important;
  text-shadow: 0 0 12px rgba(255, 255, 255, 0.6);
}

:deep(.ant-menu-item-selected a::after) {
  transform: scaleX(1);
}

:deep(.ant-menu-item-selected::after) {
  border-bottom: 2px solid #1890ff !important;
  box-shadow: 0 0 10px #1890ff;
}

/* 右侧用户区域 */
.header-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  padding-left: 20px;
}

.header-btn {
  color: rgba(255, 255, 255, 0.9) !important;
  font-weight: 500;
  padding: 0 18px;
  height: 42px;
  border-radius: 6px;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  position: relative;
  overflow: hidden;
  margin: 0 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(120deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  transform: translateX(-100%);
  transition: all 0.6s ease;
}

.header-btn:hover::before {
  transform: translateX(100%);
}

.header-btn:hover {
  background: rgba(255, 255, 255, 0.15) !important;
  color: white !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.user-btn {
  min-width: 130px;
  font-size: 15px;
}

/* 下拉菜单样式 */
:deep(.theme-menu),
:deep(.user-menu) {
  border-radius: 10px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
  padding: 6px;
  min-width: 180px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  overflow: hidden;
}

:deep(.ant-dropdown-menu-item) {
  border-radius: 6px;
  margin: 4px 0;
  padding: 10px 16px;
  transition: all 0.25s cubic-bezier(0.25, 0.8, 0.25, 1);
  display: flex;
  align-items: center;
  font-size: 14px;
}

:deep(.ant-dropdown-menu-item:hover) {
  background: rgba(24, 144, 255, 0.15);
  transform: translateX(2px);
}

:deep(.ant-dropdown-menu-item .anticon) {
  margin-right: 10px;
  font-size: 16px;
  transition: all 0.3s;
}

:deep(.ant-dropdown-menu-item:hover .anticon) {
  color: #1890ff;
  transform: scale(1.1);
}

/* 内容区域 */
.content {
  padding: 24px;
  background:rgb(84, 93, 147);
  min-height: calc(100vh - 160px);
  position: relative;
  z-index: 1;
}

/* 页面切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.4s ease, transform 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(15px);
}

/* 页脚 */
.footer {
  text-align: center;
  padding: 20px 24px;
  background: #f5f2eb;
  box-shadow: 0 -1px 4px rgba(0, 0, 0, 0.03);
  position: relative;
  z-index: 5;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
}

.copyright {
  color: rgba(0, 0, 0, 0.45);
  font-size: 14px;
}

/* 响应式调整 */
@media (max-width: 992px) {
  .header {
    padding: 0 16px;
  }
  
  .logo-text {
    font-size: 18px;
  }
  
  :deep(.ant-menu-item a) {
    font-size: 15px;
  }
  
  .header-btn {
    padding: 0 14px;
  }
}

@media (max-width: 768px) {
  .header {
    padding: 0 12px;
  }
  
  .logo {
    padding: 0 16px;
    margin-right: 10px;
  }
  
  .logo-text {
    font-size: 16px;
  }
  
  .logo-icon {
    font-size: 26px;
    margin-right: 10px;
  }
  
  :deep(.ant-menu-item) {
    padding: 0 12px;
    margin: 0 2px;
  }
  
  :deep(.ant-menu-item a) {
    font-size: 14px;
  }
  
  .header-btn {
    padding: 0 12px;
    height: 38px;
    font-size: 14px;
  }
  
  .user-btn {
    min-width: 110px;
  }
  
  .content {
    padding: 16px 12px;
  }
}
</style>