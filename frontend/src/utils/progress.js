import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

// 进度条配置
NProgress.configure({
  easing: 'ease-out',
  speed: 700,
  showSpinner: false,
  trickleSpeed: 300,
  minimum: 0.2
})

export default NProgress 