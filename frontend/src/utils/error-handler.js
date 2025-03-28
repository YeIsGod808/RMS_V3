import { message } from 'ant-design-vue'
import {
  handleHttpError,
  handleBusinessError,
  handleUnhandledRejection,
  handleRuntimeError
} from './error'

export function setupErrorHandler(app) {
  // Vue错误处理
  app.config.errorHandler = (err, vm, info) => {
    handleRuntimeError(err, vm, info)
  }
  
  // 全局未捕获错误
  window.onerror = function(msg, url, lineNo, columnNo, error) {
    console.error('全局未捕获错误:', {
      msg,
      url,
      lineNo,
      columnNo,
      error
    })
    message.error('系统错误，请刷新页面重试')
    return false
  }
  
  // Promise错误处理
  window.addEventListener('unhandledrejection', event => {
    event.preventDefault()
    handleUnhandledRejection(event)
  })
}

// 导出错误处理方法供API调用使用
export {
  handleHttpError,
  handleBusinessError
} 