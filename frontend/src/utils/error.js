import { message } from 'ant-design-vue'

// 错误码映射
const errorCodeMap = {
  400: '请求参数错误',
  401: '未登录或登录已过期',
  403: '没有权限访问',
  404: '请求的资源不存在',
  500: '服务器错误',
  502: '网关错误',
  503: '服务不可用',
  504: '网关超时'
}

// 处理HTTP错误
export function handleHttpError(error) {
  if (error.response) {
    const { status, data } = error.response
    const errorMessage = data?.message || errorCodeMap[status] || '请求失败'
    // message.error(errorMessage)
    return errorMessage
  } else if (error.request) {
    const errorMessage = '网络错误，请检查网络连接'
    message.error(errorMessage)
    return errorMessage
  } else {
    const errorMessage = error.message || '请求失败'
    message.error(errorMessage)
    return errorMessage
  }
}

// 处理业务错误
export function handleBusinessError(error) {
  const errorMessage = error.message || '操作失败'
  message.error(errorMessage)
  return errorMessage
}

// 处理未捕获的Promise错误
export function handleUnhandledRejection(event) {
  console.error('未处理的Promise错误:', event.reason)
  const errorMessage = '系统错误，请稍后重试'
  message.error(errorMessage)
  return errorMessage
}

// 处理运行时错误
export function handleRuntimeError(error, vm, info) {
  console.error('运行时错误:', {
    error,
    vm,
    info
  })
  const errorMessage = '系统错误，请刷新页面重试'
  message.error(errorMessage)
  return errorMessage
} 