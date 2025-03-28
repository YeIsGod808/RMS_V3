<template>
  <div class="video-player">
    <a-card :title="props.video?.title" size="small">
      <template #extra>
        <a :href="props.video?.url" target="_blank">
          <video-camera-outlined /> 观看视频
        </a>
      </template>
      <p>{{ props.video?.description }}</p>
    </a-card>
  </div>
</template>

<script setup>
import { VideoCameraOutlined } from '@ant-design/icons-vue'

const props = defineProps({
  video: {
    type: Object,
    required: true,
    default: () => ({
      title: '',
      url: '',
      description: ''
    })
  }
})
</script>

<style scoped>
.video-player {
  margin-bottom: 16px;
}
</style> 