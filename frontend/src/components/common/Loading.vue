<template>
  <div v-if="loading" class="loading-container">
    <a-spin size="large">
      <template #indicator>
        <loading-outlined spin />
      </template>
      <div class="loading-text">{{ text }}</div>
    </a-spin>
  </div>
</template>

<script setup>
import { LoadingOutlined } from '@ant-design/icons-vue'

defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  text: {
    type: String,
    default: '加载中...'
  }
})
</script>

<style scoped>
.loading-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.65);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.loading-text {
  margin-top: 8px;
  color: #1890ff;
}
</style> 