<template>
  <router-view v-slot="{ Component }">
    <transition name="fade" mode="out-in" :css="true" @after-leave="afterTransitionLeave">
      <div> 
        <component :is="Component" />
      </div>
    </transition>
  </router-view>
</template>

<script setup>
// 添加过渡结束时的回调函数，用于清理可能的DOM引用
function afterTransitionLeave() {
  // 此处可以添加过渡完成后的清理逻辑
  // 例如清除某些缓存或临时状态
  console.log('页面过渡动画完成')
}
</script>

<style>
#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    'Noto Sans', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* 页面切换动画 - 更平滑的过渡效果 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease, transform 0.4s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* 全局样式 */
.page-container {
  padding: 24px;
  background: #f0f2f5;
  min-height: calc(100vh - 64px);
}

.page-title {
  margin-bottom: 24px;
  color: rgba(0, 0, 0, 0.85);
  font-weight: 600;
  font-size: 20px;
  line-height: 32px;
}
</style> 