// 默认主题配置
export const defaultTheme = {
  token: {
    colorPrimary: '#1890ff',
    borderRadius: 4,
    wireframe: false
  }
}

// 暗色主题配置
export const darkTheme = {
  token: {
    colorPrimary: '#1890ff',
    borderRadius: 4,
    wireframe: false
  },
  // 使用暗色算法
  algorithm: 'dark'
}

// 紧凑主题配置
export const compactTheme = {
  token: {
    colorPrimary: '#1890ff',
    borderRadius: 2,
    wireframe: true
  },
  // 使用紧凑算法
  algorithm: 'compact'
} 