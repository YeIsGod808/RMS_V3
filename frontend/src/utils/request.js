import axios from 'axios'
import { useUserStore } from '@/store/modules/user'
import router from '@/router'
import { handleHttpError } from './error-handler'
import NProgress from './progress'

// 创建axios实例
const service = axios.create({
  baseURL: import.meta.env.VITE_APP_BASE_API || 'http://localhost:8005',
  timeout: 15000
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    NProgress.start()
    
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers['Authorization'] = `Bearer ${userStore.token}`
    }
    return config
  },
  error => {
    NProgress.done()
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    NProgress.done()
    const res = response.data
    
    // 添加调试日志
    console.log('API响应数据:', res)
    
    // 适配后端返回格式
    // 后端可能返回 ret 而不是 code
    if (res.ret !== undefined && res.ret !== "0" && res.ret !== 0) {
      const errMsg = res.msg || '请求失败'
      console.error('API请求错误 (ret):', res.ret, errMsg)
      
      handleHttpError({ 
        response: { 
          data: { 
            message: errMsg,
            code: res.ret 
          } 
        } 
      })
      
      // 处理401未授权错误
      if (res.ret === 'JWT_ERR' || res.ret === 'AUTH_ERR') {
        const userStore = useUserStore()
        userStore.clearUserInfo()
        router.push('/login')
      }
      
      return Promise.reject(new Error(errMsg))
    }
    
    // 兼容原有code格式
    if (res.code !== undefined && res.code !== 200) {
      const errMsg = res.message || '请求失败'
      console.error('API请求错误 (code):', res.code, errMsg)
      
      handleHttpError({ response })
      
      if (res.code === 401) {
        const userStore = useUserStore()
        userStore.clearUserInfo()
        router.push('/login')
      }
      
      return Promise.reject(new Error(errMsg))
    }
    
    return res
  },
  error => {
    NProgress.done()
    console.error('API请求异常:', error)
    
    handleHttpError(error)
    
    // 如果是401错误，清除用户信息并跳转到登录页
    if (error.response && error.response.status === 401) {
      const userStore = useUserStore()
      userStore.clearUserInfo()
      router.push('/login')
    }
    
    return Promise.reject(error)
  }
)

// 后端错误码
const errcode = {
  JWT_ERR: 'JWT_ERR',
  AUTH_ERR: 'AUTH_ERR',
  MISSING_PARAM: 'MISSING_PARAM',
  WRONG_PARAM: 'WRONG_PARAM',
  WRONG_PASSWORD: 'WRONG_PASSWORD',
  DB_CONN_ERR: 'DB_CONN_ERR',
  DB_DUP_ERR: 'DB_DUP_ERR'
}

export default service 