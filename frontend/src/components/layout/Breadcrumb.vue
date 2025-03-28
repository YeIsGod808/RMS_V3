<template>
  <div class="breadcrumb-container">
    <a-breadcrumb>
      <a-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">
        <a @click="handleClick(item)">{{ item.name }}</a>
      </a-breadcrumb-item>
    </a-breadcrumb>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useGraphStore } from '@/store/modules/graph'

const graphStore = useGraphStore()

const breadcrumbs = computed(() => {
  const { currentChapter, currentSection } = graphStore
  const crumbs = [{ name: '知识图谱', path: 'root' }]
  
  if (currentChapter) {
    crumbs.push({
      name: currentChapter.name,
      path: `chapter-${currentChapter.id}`
    })
  }
  
  if (currentSection) {
    crumbs.push({
      name: currentSection.name,
      path: `section-${currentSection.id}`
    })
  }
  
  return crumbs
})

const handleClick = (item) => {
  if (item.path === 'root') {
    graphStore.loadChapterGraph()
  } else if (item.path.startsWith('chapter-')) {
    const chapterId = item.path.split('-')[1]
    graphStore.loadSectionGraph(chapterId)
  }
}
</script>

<style scoped>
.breadcrumb-container {
  margin-bottom: 16px;
}

:deep(a) {
  color: #1890ff;
  cursor: pointer;
}

:deep(a:hover) {
  color: #40a9ff;
}
</style> 