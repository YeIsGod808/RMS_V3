import { defineStore } from 'pinia'
import { login, logout, checkToken } from '@/api/user'
import { message } from 'ant-design-vue'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    username: localStorage.getItem('username') || '',
    role: localStorage.getItem('role') || '',
    nickname: localStorage.getItem('nickname') || '',
    currentUserId: '',
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.role === 'admin',
    isTeacher: (state) => state.role === 'teacher',
    isStudent: (state) => state.role === 'student',
  },
  
  actions: {
    async login(username, password) {
      try {
        const response = await login(username, password)
        if (response.ret === "0") {
          // 后端返回用户信息，包含token
          const userData = response.user
          this.token = userData.token
          this.username = userData.id
          this.nickname = userData.nickname
          this.role = userData.usertype
          
          // 保存到本地存储
          localStorage.setItem('token', userData.token)
          localStorage.setItem('username', userData.id)
          localStorage.setItem('nickname', userData.nickname)
          localStorage.setItem('role', userData.usertype)
          
          return userData
        } else {
          // 登录失败，但有返回信息
          throw new Error(response.msg || '登录失败')
        }
      } catch (error) {
        message.error('登录失败')
        throw error
      }
    },
    
    async verifyToken() {
      if (!this.token) return false
      
      try {
        const response = await checkToken(this.token)
        if (response.ret === "0") {
          // 更新用户信息
          const userData = response.user
          this.username = userData.id
          this.nickname = userData.nickname
          this.role = userData.usertype
          
          // 更新本地存储
          localStorage.setItem('username', userData.id)
          localStorage.setItem('nickname', userData.nickname)
          localStorage.setItem('role', userData.usertype)
          
          return true
        } else {
          this.clearUserInfo()
          return false
        }
      } catch (error) {
        console.error('Token验证失败:', error)
        this.clearUserInfo()
        return false
      }
    },
    
    async logout() {
      try {
        // 尝试调用后端登出API，但不依赖其结果
        await logout()
      } catch (error) {
        console.error('登出API调用失败:', error)
      }
      
      // 无论API调用成功与否，都清除用户信息
      this.clearUserInfo()
    },

    clearUserInfo() {
      // 清除状态
      this.token = ''
      this.username = ''
      this.role = ''
      this.nickname = ''
      
      // 清除本地存储
      localStorage.removeItem('token')
      localStorage.removeItem('username')
      localStorage.removeItem('role')
      localStorage.removeItem('nickname')
    },
    
    // 设置当前操作的用户ID
    setCurrentUserId(userId) {
      this.currentUserId = userId
    }
  }
}) 