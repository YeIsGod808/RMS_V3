import { defineStore } from 'pinia'
import { defaultTheme, darkTheme, compactTheme } from '@/utils/theme'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    theme: 'default',
    themeConfig: defaultTheme,
    sidebarCollapsed: false
  }),

  actions: {
    setTheme(theme) {
      console.log('切换主题:', theme)
      this.theme = theme
      
      // 设置主题配置
      switch (theme) {
        case 'dark':
          this.themeConfig = darkTheme
          document.body.setAttribute('data-theme', 'dark')
          break
        case 'compact':
          this.themeConfig = compactTheme
          document.body.setAttribute('data-theme', 'compact')
          break
        default:
          this.themeConfig = defaultTheme
          document.body.setAttribute('data-theme', 'default')
      }
      
      // 保存到本地存储
      localStorage.setItem('theme', theme)
      
      // 通知主题变化，强制重新渲染
      this.$patch({ themeUpdated: Date.now() })
      
      // 应用CSS变量
      this.applyThemeVariables()
    },

    applyThemeVariables() {
      // 获取当前使用的主题配置
      const currentTheme = this.themeConfig
      
      // 将主题配置中的token应用为CSS变量
      if (currentTheme && currentTheme.token) {
        Object.entries(currentTheme.token).forEach(([key, value]) => {
          // 转换驼峰命名为CSS变量命名
          const cssVarName = key.replace(/([A-Z])/g, '-$1').toLowerCase()
          document.documentElement.style.setProperty(`--${cssVarName}`, value)
        })
      }
    },

    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed
      // 保存到本地存储
      localStorage.setItem('sidebarCollapsed', String(this.sidebarCollapsed))
    },

    initSettings() {
      // 从本地存储加载设置
      const theme = localStorage.getItem('theme')
      if (theme) {
        this.setTheme(theme)
      } else {
        // 确保默认主题也应用CSS变量
        this.applyThemeVariables()
      }
      
      const sidebarCollapsed = localStorage.getItem('sidebarCollapsed')
      if (sidebarCollapsed !== null) {
        this.sidebarCollapsed = sidebarCollapsed === 'true'
      }
    }
  }
}) 