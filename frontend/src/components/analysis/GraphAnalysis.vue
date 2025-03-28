<template>
  <div class="analysis-container">
    <a-row :gutter="16">
      <a-col :span="8">
        <a-card title="知识点关联度分析">
          <ve-bar :data="connectionData" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card title="学习难度分布">
          <ve-pie :data="difficultyData" />
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card title="学习路径推荐">
          <learning-path-chart :path="recommendedPath" />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useGraphStore } from '@/store/modules/graph'
import { analyzeConnections, getDifficultyDistribution } from '@/api/analysis'

const graphStore = useGraphStore()
const connectionData = ref([])
const difficultyData = ref([])
const recommendedPath = ref([])

onMounted(async () => {
  const analysisData = await analyzeConnections()
  connectionData.value = formatConnectionData(analysisData)
  
  const difficultyStats = await getDifficultyDistribution()
  difficultyData.value = formatDifficultyData(difficultyStats)
})
</script>