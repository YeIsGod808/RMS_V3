<template>
  <div class="error-container">
    <div class="error-content">
      <h1 class="error-code">403</h1>
      <h2 class="error-title">禁止访问</h2>
      <p class="error-desc">抱歉，您没有权限访问该页面</p>
      <a-button type="primary" @click="goBack">返回上一页</a-button>
      <a-button @click="goHome" style="margin-left: 16px;">返回首页</a-button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => {
  router.go(-1)
}

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
.error-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: var(--background-color, #f0f2f5);
}

.error-content {
  text-align: center;
  padding: 48px 32px;
  background-color: var(--component-background, #fff);
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.error-code {
  font-size: 72px;
  font-weight: 600;
  color: #f5222d;
  margin-bottom: 8px;
}

.error-title {
  font-size: 24px;
  font-weight: 500;
  color: var(--heading-color, rgba(0, 0, 0, 0.85));
  margin-bottom: 16px;
}

.error-desc {
  font-size: 16px;
  color: var(--text-color-secondary, rgba(0, 0, 0, 0.45));
  margin-bottom: 24px;
}
</style> 
 