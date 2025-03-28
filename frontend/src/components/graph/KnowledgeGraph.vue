<template>
  <div class="knowledge-graph">
    <Loading :loading="graphStore.loading" />
    
    <div class="graph-toolbar">
      <a-space>
        <a-input-search
          v-model:value="searchText"
          placeholder="搜索节点"
          @search="onSearch"
          style="width: 200px"
        />
        <a-button-group>
          <a-tooltip title="放大">
            <a-button @click="zoomIn"><plus-outlined /></a-button>
          </a-tooltip>
          <a-tooltip title="缩小">
            <a-button @click="zoomOut"><minus-outlined /></a-button>
          </a-tooltip>
          <a-tooltip title="重置缩放">
            <a-button @click="resetZoom"><redo-outlined /></a-button>
          </a-tooltip>
          <a-tooltip title="重新排列">
            <a-button @click="reloadGraph"><reload-outlined /></a-button>
          </a-tooltip>
          <a-tooltip title="返回上一级" v-if="graphStore.currentLevel !== 'chapter'">
            <a-button @click="goBackToParent"><rollback-outlined /></a-button>
          </a-tooltip>
        </a-button-group>
      </a-space>
    </div>
    
    <!-- 添加图例说明 -->
    <div class="legend-panel">
      <!-- <h3>学习状态图例</h3>
      <div class="legend-items">
        <div class="legend-item">
          <span class="legend-color" style="background-color: #E74C3C;"></span>
          <span class="legend-label">未学习</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background-color: #F39C12;"></span>
          <span class="legend-label">正在学习</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background-color: #2ECC71;"></span>
          <span class="legend-label">已完成</span>
        </div>
      </div> -->
      
      <h3>节点类型</h3>
      <div class="legend-items">
        <div class="legend-item">
          <span class="legend-color" style="background-color: #1143D7;"></span>
          <span class="legend-label">章节</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background-color: #722ed1;"></span>
          <span class="legend-label">小节</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background-color: #52C41A;"></span>
          <span class="legend-label">知识点</span>
        </div>
      </div>
      
      <h3>关系类型</h3>
      <div class="legend-items">
        <div class="legend-item">
          <span class="legend-color" style="background-color: #FA8C16;"></span>
          <span class="legend-label">前置</span>
          <span class="legend-text">学习顺序依赖</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background-color: #1890FF;"></span>
          <span class="legend-label">包含</span>
          <span class="legend-text">层级包含关系</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background-color: #52C41A;"></span>
          <span class="legend-label">相关</span>
          <span class="legend-text">相关知识</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background-color: #722ED1;"></span>
          <span class="legend-label">扩展</span>
          <span class="legend-text">扩展内容</span>
        </div>
      </div>
    </div>
    
    <div ref="graphContainer" class="graph-container">
      <!-- 空数据提示 -->
      <div v-if="!graphStore.loading && (!graphStore.currentGraphData || !graphStore.currentGraphData.nodes || !graphStore.currentGraphData.nodes.length)" class="empty-data">
        <a-empty description="暂无图谱数据">
          <template #extra>
            <a-button type="primary" @click="reloadGraph">
              重新加载
            </a-button>
          </template>
        </a-empty>
      </div>
      
      <!-- 错误提示 -->
      <div v-if="hasError" class="error-message">
        <div class="title">加载图谱时出现错误</div>
        <div class="content">{{ errorMessage }}</div>
        <a-button type="primary" @click="reloadGraph" style="margin-top: 16px">
          重新加载
        </a-button>
      </div>
    </div>
    
    <!-- 搜索结果抽屉 -->
    <a-drawer
      v-model:visible="searchResultsVisible"
      title="搜索结果"
      placement="right"
      width="400"
    >
      <template v-if="searchResults && searchResults.length">
        <a-list :data-source="searchResults" size="small">
          <template #header>
            <div>找到 {{ searchResults.length }} 个匹配的节点</div>
          </template>
          <template #renderItem="{ item }">
            <a-list-item>
              <a-list-item-meta>
                <template #title>
                  <a @click="navigateToNode(item)">{{ item.name }}</a>
                </template>
                <template #description>
                  <div>
                    <a-tag :color="getNodeTypeColor(item.type)">{{ getNodeTypeName(item.type) }}</a-tag>
                    <span style="margin-left: 8px">{{ item.description ? (item.description.slice(0, 30) + (item.description.length > 30 ? '...' : '')) : '' }}</span>
                  </div>
                </template>
              </a-list-item-meta>
            </a-list-item>
          </template>
        </a-list>
      </template>
      <template v-else>
        <a-empty description="未找到匹配的节点" />
      </template>
    </a-drawer>
    
    <!-- 知识点详情抽屉 -->
    <a-drawer
      v-model:visible="drawerVisible"
      title="知识点详情"
      placement="right"
      :width="500"
    >
      <div v-if="selectedNode" class="knowledge-detail">
        <h2 class="knowledge-title">{{ selectedNode.name }}</h2>
        
        <a-card class="detail-card">
          <template #title><h3 class="section-title">描述</h3></template>
          <div class="description-content">
            {{ selectedNode.description || '暂无描述' }}
          </div>
        </a-card>
        <a-card class="detail-card">
          <template #title><h3 class="section-title">学习难度评估</h3></template>
          <div class="difficulty-assessment">
            <div class="difficulty-score">
              <div class="score-label">难度评分</div>
              <a-rate v-model:value="difficultyScore" :count="5" disabled allow-half />
              <div class="score-value">{{ difficultyScore }}/5</div>
            </div>
            <a-progress
              :percent="(difficultyScore / 5) * 100"
              :stroke-color="{
                '0%': '#87d068',
                '100%': '#ff4d4f'
              }"
              :show-info="false"
              class="difficulty-progress"
            />
            <div class="difficulty-stats">
              <div class="stat-item">
                <div class="stat-label">前置知识点</div>
                <div class="stat-value">
                  <number-outlined />
                  <span>{{ prereqCount }}个</span>
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-label">最大深度</div>
                <div class="stat-value">
                  <node-index-outlined />
                  <span>{{ maxDepth }}层</span>
                </div>
              </div>
            </div>
          </div>
        </a-card>
        
        <a-card class="detail-card" v-if="selectedNode.videos && selectedNode.videos.length > 0">
          <template #title>
            <h3 class="section-title">
              <video-camera-outlined /> 相关视频
              <span class="resource-count">({{ selectedNode.videos.length }})</span>
            </h3>
          </template>
          <a-list size="small" :data-source="selectedNode.videos">
            <template #renderItem="{ item }">
              <a-list-item class="resource-list-item">
                <a-list-item-meta>
                  <template #title>
                    <div class="resource-title">
                      <a :href="item.url" target="_blank">{{ item.title }}</a>
                    </div>
                  </template>
                  <template #description>
                    <div class="resource-description">{{ item.description || '暂无描述' }}</div>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-button type="primary" size="small" @click="previewVideo(item)">
                    预览
                  </a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
        
        <a-card class="detail-card" v-if="selectedNode.courseware && selectedNode.courseware.length > 0">
          <template #title>
            <h3 class="section-title">
              <book-outlined /> 相关课件
              <span class="resource-count">({{ selectedNode.courseware.length }})</span>
            </h3>
          </template>
          <a-list size="small" :data-source="selectedNode.courseware">
            <template #renderItem="{ item }">
              <a-list-item class="resource-list-item">
                <a-list-item-meta>
                  <template #title>
                    <div class="resource-title">
                      <a :href="item.url" target="_blank">{{ item.title }}</a>
                    </div>
                  </template>
                  <template #description>
                    <div class="resource-description">{{ item.description || '暂无描述' }}</div>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-button type="primary" size="small" @click="previewCourseware(item)">
                    预览
                  </a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
        
        <a-card class="detail-card" v-if="selectedNode.exercises && selectedNode.exercises.length > 0">
          <template #title>
            <h3 class="section-title">
              <file-outlined /> 相关练习题
              <span class="resource-count">({{ selectedNode.exercises.length }})</span>
            </h3>
          </template>
          <a-list size="small" :data-source="selectedNode.exercises">
            <template #renderItem="{ item }">
              <a-list-item class="resource-list-item">
                <a-list-item-meta>
                  <template #title>
                    <div class="resource-title">
                      <a :href="item.url" target="_blank">{{ item.title }}</a>
                    </div>
                  </template>
                  <template #description>
                    <div class="resource-description">
                      <a-tag 
                        :color="getDifficultyColor(item.difficulty)" 
                        style="margin-right: 8px;"
                      >
                        {{ getDifficultyText(item.difficulty) }}
                      </a-tag>
                      <span>{{ item.description || '暂无描述' }}</span>
                    </div>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-button type="primary" size="small" @click="previewExercise(item)">
                    预览
                  </a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
        
        <div class="no-resources" v-if="(!selectedNode.videos || selectedNode.videos.length === 0) &&
                                      (!selectedNode.courseware || selectedNode.courseware.length === 0) && 
                                      (!selectedNode.exercises || selectedNode.exercises.length === 0)">
          <a-empty description="暂无相关资源" />
        </div>
      </div>
      <div v-else class="empty-data">
        <a-empty description="未找到知识点详情" />
      </div>
    </a-drawer>
    
    <!-- 预览模态框 -->
    <a-modal
      v-model:visible="previewModalVisible"
      :title="previewModalTitle"
      :footer="null"
      :width="800"
      @cancel="handlePreviewModalClose"
    >
      <div v-if="previewContent && previewContent.url">
        <template v-if="previewContentType === 'video' && previewContent.cover_url">
          <video 
            controls 
            style="width: 100%; max-height: 500px;" 
            :src="previewContent.url"
            :poster="previewContent.cover_url"
            ref="videoPlayer"
          ></video>
        </template>
        <template v-else-if="previewContentType === 'pdf'">
          <div style="height: 500px; width: 100%;">
            <iframe 
              :src="previewContent.url" 
              style="width: 100%; height: 100%; border: none;"
            ></iframe>
          </div>
        </template>
        <template v-else>
          <div style="text-align: center; padding: 20px;">
            <p>无法直接预览该资源，请点击下方按钮跳转查看</p>
            <a :href="previewContent.url" target="_blank">
              <a-button type="primary">跳转查看</a-button>
            </a>
          </div>
        </template>
      </div>
      <div v-else style="text-align: center; padding: 20px;">
        <a-empty description="资源URL不存在" />
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick, onUnmounted, reactive } from 'vue'
import * as d3 from 'd3'
import { useGraphStore } from '@/store/modules/graph'
import { message } from 'ant-design-vue'
import { getDifficultyDistribution } from '@/api/analysis'
import { 
  PlusOutlined, 
  MinusOutlined, 
  RedoOutlined,
  ReloadOutlined,
  RollbackOutlined,
  VideoCameraOutlined,
  BookOutlined,
  NodeIndexOutlined,
  NumberOutlined,
  FileOutlined
} from '@ant-design/icons-vue'
import Loading from '../common/Loading.vue'

const graphStore = useGraphStore()
const graphContainer = ref(null)
const searchText = ref('')
const drawerVisible = ref(false)
const selectedNode = ref(null)
const searchResults = ref([])
const searchResultsVisible = ref(false)
const difficultyScore = ref(0)
const prereqCount = ref(0)
const maxDepth = ref(0)

// 获取学习难度评估数据
const fetchDifficultyData = async (nodeId) => {
  try {
    const res = await getDifficultyDistribution(nodeId)
    if (res.code === 200) {
      // 将0-5的难度分数转换为0-5的评分
      difficultyScore.value = res.data.difficulty_score
      prereqCount.value = res.data.prereq_count
      maxDepth.value = res.data.max_depth
    }
  } catch (error) {
    console.error('获取学习难度数据失败:', error)
  }
}

const mockData = ref({
  sections: [],
  points: []
})

let svg = null
let g = null
let zoom = null
let resizeObserver = null
let resizeTimer = null

const hasError = ref(false)
const errorMessage = ref('')

// 图谱配置
const graphConfig = {
  width: 800,
  height: 600,
  nodeRadius: {
    chapter: 50,
    section: 40,
    knowledge: 30
  },
  nodeColor: {
    chapter: 'rgb(5, 32, 111)',
    section: '#52c41a',
    knowledge: '#722ed1'
  },
  linkDistance: 300,
  linkStrength: 0.3,
  chargeStrength: -2500,
  collideStrength: 1.2,
  centerStrength: 0.03
}

// 调整容器尺寸时更新图谱
const handleResize = () => {
  if (!graphContainer.value || !svg) return
  
  // 使用防抖处理，避免频繁触发
  if (resizeTimer) clearTimeout(resizeTimer)
  
  resizeTimer = setTimeout(() => {
    const container = graphContainer.value
    const width = container.clientWidth
    const height = container.clientHeight
    
    svg.attr('width', width)
      .attr('height', height)
    
    // 如果有当前数据，重新渲染图谱
    if (graphStore.currentGraphData) {
      console.log('容器尺寸变化，重新渲染图谱')
      updateGraph()
    }
  }, 200) // 200ms防抖
}

// 初始化图谱
const initGraph = () => {
  const container = graphContainer.value
  const width = container.clientWidth
  const height = container.clientHeight
  
  // 创建SVG
  svg = d3.select(container)
    .append('svg')
    .attr('width', width)
    .attr('height', height)
  
  // 创建缩放行为
  zoom = d3.zoom()
    .scaleExtent([0.1, 4])
    .on('zoom', (event) => {
      g.attr('transform', event.transform)
    })
  
  svg.call(zoom)
  
  // 创建容器组
  g = svg.append('g')
  
  // 监听图谱数据变化
  watch(() => graphStore.currentGraphData, (data) => {
    if (data) {
      updateGraph()
    }
  }, { deep: true })
}

// 组件挂载时初始化
onMounted(async () => {
  try {
    console.log('组件挂载，初始化图谱')
    await nextTick() // 确保DOM已渲染
    
    if (!graphContainer.value) {
      console.error('图谱容器未找到')
      hasError.value = true
      errorMessage.value = '图谱容器初始化失败'
      return
    }
    
    // 设置监听器以处理容器尺寸变化
    resizeObserver = new ResizeObserver(handleResize)
    resizeObserver.observe(graphContainer.value)
    
    initGraph()
    
    // 初始化模拟数据结构
    initMockData()
    
    // 打印调试信息
    console.log('当前图谱数据状态:', {
      hasData: !!graphStore.currentGraphData,
      loading: graphStore.loading,
      level: graphStore.currentLevel
    })
    
    // 如果没有当前图谱数据，加载章节图谱
    if (!graphStore.currentGraphData) {
      console.log('加载初始章节图谱')
      await graphStore.loadChapterGraph()
    } else {
      // 已有数据，直接更新图谱
      console.log('使用现有数据更新图谱')
      updateGraph()
    }
    
    // 重置错误状态
    hasError.value = false
    errorMessage.value = ''
  } catch (error) {
    console.error('初始化图谱错误:', error)
    hasError.value = true
    errorMessage.value = '初始化图谱失败，请刷新页面重试'
    // message.error('初始化图谱失败，请刷新页面重试')
  }
})

// 组件卸载时清理
onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
  
  if (resizeTimer) {
    clearTimeout(resizeTimer)
  }
})

// 获取节点数据，不再进行筛选
const getGraphNodes = () => {
  if (!graphStore.currentGraphData) return []
  return graphStore.currentGraphData.nodes || []
}

// 更新图谱
const updateGraph = () => {
  // 清除错误状态
  hasError.value = false
  errorMessage.value = ''
  
  // 获取当前图谱数据
  const data = graphStore.currentGraphData
  if (!data || !data.nodes || data.nodes.length === 0) {
    console.log('无图谱数据')
    return
  }
  
  // 检查图谱容器是否存在
  if (!graphContainer.value) {
    console.error('图谱容器不存在')
    hasError.value = true
    errorMessage.value = '图谱容器未初始化，请刷新页面重试'
    return
  }
  
  console.log(`渲染图谱: ${data.nodes.length}个节点, ${data.links ? data.links.length : 0}个连接`)
  console.log('图谱容器尺寸:', graphContainer.value.clientWidth, graphContainer.value.clientHeight)
  
  // 应用筛选器 - 改为直接获取节点
  const nodes = JSON.parse(JSON.stringify(getGraphNodes()))
  const links = data.links ? JSON.parse(JSON.stringify(data.links)) : []
  
  // 处理多个关系 - 给每个关系创建唯一标识
  links.forEach((link, index) => {
    link.id = link.id || `link-${index}`
    
    // 检查是否有重复的source-target对
    const sameLinkCount = links.filter((l, i) => 
      i < index && 
      ((l.source === link.source && l.target === link.target) || 
       (l.target === link.source && l.source === link.target))
    ).length
    
    // 为重复链接添加位移信息
    if (sameLinkCount > 0) {
      link.duplicate = true
      link.duplicateIndex = sameLinkCount
    }
  })
  
  // 清空原有图谱
  if (g) {
    g.selectAll('*').remove()
  } else {
    svg.selectAll('*').remove()
    g = svg.append('g')
  }
  
  // 创建力导向图
  const simulation = d3.forceSimulation(nodes)
    .force('link', d3.forceLink(links).id(d => d.id).distance(120))
    .force('charge', d3.forceManyBody().strength(-400))
    .force('center', d3.forceCenter(graphContainer.value.clientWidth / 2, graphContainer.value.clientHeight / 2))
  
  // 定义连接线类型和颜色映射
  const linkColors = {
    '前置': '#FA8C16', // 前置关系 - 橙色
    '包含': '#1890FF', // 包含关系 - 蓝色
    '相关': '#52C41A', // 相关关系 - 绿色
    '扩展': '#722ED1'  // 扩展关系 - 紫色
  }
  
  // 定义箭头标记
  const defs = svg.append('defs')
  
  const arrowTypes = Object.keys(linkColors)
  arrowTypes.forEach(type => {
    defs.append('marker')
      .attr('id', `arrow-${type}`)
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', 0)  
      .attr('refY', 0)
      .attr('markerWidth', 5)
      .attr('markerHeight', 5)
      .attr('orient', 'auto')
      .append('path')
      .attr('d', 'M0,-5L10,0L0,5')
      .attr('fill', linkColors[type])
  })
  
  // 创建路径生成器用于绘制曲线
  const linkPath = d3.linkHorizontal()
    .x(d => d.x)
    .y(d => d.y);
  
  // 绘制连接线 - 使用路径替代直线
  const link = g.append('g')
    .attr('class', 'links')
    .selectAll('path')
    .data(links)
    .enter().append('path')
    .attr('stroke-width', 2.5)
    .attr('fill', 'none')
    .attr('stroke', d => linkColors[d.type] || '#999')
    .attr('stroke-dasharray', d => d.type === '相关' ? '5,5' : null)  // 相关关系使用虚线
    .attr('marker-end', d => `url(#arrow-${d.type})`)
  
  // 添加连接线标签
  const linkLabels = g.append('g')
    .attr('class', 'link-labels')
    .selectAll('g')
    .data(links)
    .enter().append('g')
    .attr('class', 'link-label')
  
  // 添加标签背景
  linkLabels.append('rect')
    .attr('rx', 4)
    .attr('ry', 4)
    .attr('fill', 'white')
    .attr('opacity', 0) 
    .attr('width', 35)  // 默认宽度，将在tick函数中动态调整
    .attr('height', 20)
  
  // 添加标签文本
  linkLabels.append('text')
    .attr('text-anchor', 'middle')
    .attr('dominant-baseline', 'middle')
    .attr('font-size', 12)
    .attr('fill', d => linkColors[d.type] || '#666')
    .text(d => d.type || '关系')
  
  // 绘制节点
  const node = g.append('g')
    .attr('class', 'nodes')
    .selectAll('.node')
    .data(nodes)
    .enter().append('g')
    .attr('class', 'node')
    .on('click', (event, d) => handleNodeClick(event, d))
    .call(d3.drag()
      .on('start', dragstarted)
      .on('drag', dragged)
      .on('end', dragended))
  
  // 节点圆形
  node.append('circle')
    .attr('r', d => {
      // 根据节点类型和文本长度设置不同大小
      if (d.type === 'chapter') return 40;
      if (d.type === 'section') return 35;
      return 30; // 知识点默认大小
    })
    .style('fill', d => {
      // 根据节点类型和学习状态设置颜色
      if (d.type === 'chapter') return '#1890ff'; // 章节蓝色
      if (d.type === 'section') return '#722ed1'; // 小节紫色
      
      // 知识点根据学习状态设置颜色
      const status = getLearningStatus(d.id);
      if (status === 'completed') return '#2ECC71'; // 绿色表示已完成
      if (status === 'in_progress') return '#F39C12'; // 橙色表示进行中
      if (status === 'not_started') return '#E74C3C'; // 红色表示未开始
      return '#52C41A'; // 默认绿色
    })
    .style('stroke', '#fff')
    .style('stroke-width', 2)
  
  // 节点标签修改为放在节点内部
  node.append('text')
    .attr('text-anchor', 'middle')
    .attr('dominant-baseline', 'middle')
    .text(d => d.name)
    .style('font-size', d => {
      if (d.type === 'chapter') return '14px';
      if (d.type === 'section') return '12px';
      return '10px'; // 知识点字体更小
    })
    .style('fill', '#fff')
    .style('font-weight', 'bold')
    .style('pointer-events', 'none')
    .each(function(d) {
      // 文本自动换行处理
      const text = d3.select(this);
      const words = d.name.split(/\s+/).reverse();
      const lineHeight = 1.1; // 行高
      const radius = d.type === 'chapter' ? 40 : (d.type === 'section' ? 35 : 30);
      const maxWidth = radius * 1.5; // 限制文本宽度
      
      let word;
      let line = [];
      let lineNumber = 0;
      let tspan = text.text(null).append("tspan")
        .attr("x", 0)
        .attr("y", 0)
        .attr("dy", 0);
      
      // 如果文本很短，直接显示
      if (d.name.length <= 4) {
        tspan.text(d.name);
      } else {
        // 否则按字符分行显示（简化处理，中文环境下）
        const chars = d.name.split('');
        const maxCharsPerLine = d.type === 'knowledge' ? 4 : 5;
        
        for (let i = 0; i < chars.length; i += maxCharsPerLine) {
          if (i > 0) {
            tspan = text.append("tspan")
              .attr("x", 0)
              .attr("dy", lineHeight + "em");
          }
          tspan.text(chars.slice(i, i + maxCharsPerLine).join(''));
        }
      }
    });
  
  // 更新力导向图模拟
  simulation.on('tick', () => {
    // 更新连接线路径 - 使用曲线替代直线
    link.attr('d', d => {
      // 获取节点半径，基于节点类型
      const sourceRadius = d.source.type === 'chapter' ? 40 : 
                          (d.source.type === 'section' ? 35 : 30);
      const targetRadius = d.target.type === 'chapter' ? 40 : 
                          (d.target.type === 'section' ? 35 : 30);
      
      // 计算源节点和目标节点之间的距离和角度
      const dx = d.target.x - d.source.x;
      const dy = d.target.y - d.source.y;
      const dr = Math.sqrt(dx * dx + dy * dy);
      const angle = Math.atan2(dy, dx);
      
      // 计算源节点的边缘点
      const sourceX = d.source.x + sourceRadius * Math.cos(angle);
      const sourceY = d.source.y + sourceRadius * Math.sin(angle);
      
      // 为箭头保留空间，确保不被目标节点遮挡
      // 箭头大小为5，添加5像素额外空间作为缓冲
      const arrowSpace = 13;
      
      // 计算目标节点的边缘点，留出箭头空间
      const targetX = d.target.x - (targetRadius + arrowSpace) * Math.cos(angle);
      const targetY = d.target.y - (targetRadius + arrowSpace) * Math.sin(angle);
      
      // 为重复的连接设置不同的曲率
      if (d.duplicate) {
        // 根据重复索引增加曲率
        const curvature = 0.2 + (d.duplicateIndex * 0.1);
        
        // 计算曲线的控制点
        const midX = (sourceX + targetX) / 2;
        const midY = (sourceY + targetY) / 2;
        
        // 创建控制点的偏移 - 垂直于主轴线
        const offsetX = -dy * curvature; 
        const offsetY = dx * curvature;
        
        // 控制点
        const cpx = midX + offsetX;
        const cpy = midY + offsetY;
        
        // 生成路径
        return `M${sourceX},${sourceY} Q${cpx},${cpy} ${targetX},${targetY}`;
      } else {
        // 非重复的链接使用直线
        return `M${sourceX},${sourceY} L${targetX},${targetY}`;
      }
    });
    
    // 更新连接线标签位置 - 改进标签定位算法避免遮挡箭头
    linkLabels.attr('transform', d => {
      // 如果是曲线连接（重复连接）
      if (d.duplicate) {
        // 计算源节点和目标节点之间的连线
        const dx = d.target.x - d.source.x;
        const dy = d.target.y - d.source.y;
        const angle = Math.atan2(dy, dx);
        
        // 计算中点
        const midX = (d.source.x + d.target.x) / 2;
        const midY = (d.source.y + d.target.y) / 2;
        
        // 根据重复索引计算曲率
        const curvature = 0.2 + (d.duplicateIndex * 0.1);
        
        // 垂直偏移
        const offsetX = -dy * curvature; 
        const offsetY = dx * curvature;
        
        // 标签位置 - 放在曲线的中点位置
        return `translate(${midX + offsetX}, ${midY + offsetY})`;
      } else {
        // 非重复连接的标签位置
        const midX = (d.source.x + d.target.x) / 2;
        const midY = (d.source.y + d.target.y) / 2;
        
        // 计算连线的角度
        const dx = d.target.x - d.source.x;
        const dy = d.target.y - d.source.y;
        const angle = Math.atan2(dy, dx);
        
        // 将标签上移一定距离，垂直于连线方向
        const offset = 12;
        
        const labelX = midX - offset * Math.sin(angle);
        const labelY = midY + offset * Math.cos(angle);
        
        return `translate(${labelX}, ${labelY})`;
      }
    });
    
    // 动态调整标签背景大小和位置 - 添加额外的填充以确保文本完全包含
    linkLabels.select('rect')
      .attr('width', function() {
        return d3.select(this.parentNode).select('text').node().getBBox().width + 12;
      })
      .attr('height', function() {
        return d3.select(this.parentNode).select('text').node().getBBox().height + 8;
      })
      .attr('x', function() {
        return -d3.select(this.parentNode).select('text').node().getBBox().width / 2 - 6;
      })
      .attr('y', function() {
        return -d3.select(this.parentNode).select('text').node().getBBox().height / 2 - 4;
      });
    
    // 更新节点位置
    node
      .attr('transform', d => `translate(${d.x},${d.y})`)
  })
  
  // 处理拖拽事件
  function dragstarted(event, d) {
    if (!event.active) simulation.alphaTarget(0.3).restart()
    d.fx = d.x
    d.fy = d.y
  }
  
  function dragged(event, d) {
    d.fx = event.x
    d.fy = event.y
  }
  
  function dragended(event, d) {
    if (!event.active) simulation.alphaTarget(0)
    d.fx = null
    d.fy = null
  }
}

// 工具函数
const getNodeTypeName = (type) => {
  const types = {
    chapter: '章节',
    section: '小节',
    point: '知识点'
  }
  return types[type] || type
}

// 根据难度返回颜色
const getDifficultyColor = (difficulty) => {
  switch (difficulty) {
    case 'easy': return 'green';
    case 'medium': return 'blue';
    case 'hard': return 'red';
    default: return 'default';
  }
}

// 根据难度返回文本
const getDifficultyText = (difficulty) => {
  switch (difficulty) {
    case 'easy': return '简单';
    case 'medium': return '中等';
    case 'hard': return '困难';
    default: return '未知';
  }
}

// 交互方法
const zoomIn = () => {
  svg.transition().call(zoom.scaleBy, 1.5)
}

const zoomOut = () => {
  svg.transition().call(zoom.scaleBy, 0.67)
}

const resetZoom = () => {
  svg.transition().call(zoom.transform, d3.zoomIdentity)
}

// 重新加载图谱
const reloadGraph = () => {
  try {
    hasError.value = false
    errorMessage.value = ''
    
    const currentLevel = graphStore.currentLevel
    if (currentLevel === 'chapter') {
      graphStore.loadChapterGraph()
    } else if (currentLevel === 'section' && graphStore.currentChapter) {
      graphStore.loadSectionGraph(graphStore.currentChapter.id)
    } else if (currentLevel === 'point' && graphStore.currentSection) {
      graphStore.loadPointGraph(graphStore.currentSection.id)
    }
  } catch (error) {
    console.error('重新加载图谱错误:', error)
    hasError.value = true
    errorMessage.value = '重新加载图谱失败'
    message.error('重新加载图谱失败')
  }
}

// 节点点击事件
const handleNodeClick = async (event, node) => {
  if (node) {
    fetchDifficultyData(node.id)
  }
  try {
    // 阻止事件冒泡
    event.stopPropagation();
    
    if (!node) {
      console.warn('无效的节点数据:', node);
      return;
    }
    
    console.log('Node clicked:', node);
    
    if (node.type === 'chapter') {
      console.log('加载章节小节图谱:', node.id);
      await graphStore.loadSectionGraph(node.id);
      // 更新mockData中的sections
      if (graphStore.currentGraphData && graphStore.currentGraphData.nodes) {
        mockData.value.sections = graphStore.currentGraphData.nodes.filter(n => n.type === 'section');
        console.log('更新sections数据:', mockData.value.sections);
      }
    } else if (node.type === 'section') {
      console.log('加载小节知识点图谱:', node.id);
      await graphStore.loadPointGraph(node.id);
      // 更新mockData中的points
      if (graphStore.currentGraphData && graphStore.currentGraphData.nodes) {
        mockData.value.points = graphStore.currentGraphData.nodes.filter(n => n.type === 'knowledge' || n.type === 'point');
        console.log('更新points数据:', mockData.value.points);
      }
    } else if (node.type === 'knowledge' || node.type === 'point') {
      console.log('加载知识点详情，ID:', node.id);
      message.loading({ content: '加载知识点详情...', key: 'pointDetails' });
      
      try {
        // 调用store方法加载详情
        const detailsResult = await graphStore.loadPointDetails(node.id);
        
        if (detailsResult) {
          selectedNode.value = detailsResult;
          console.log('知识点详情加载成功:', selectedNode.value);
          
          // 检查资源是否存在
          const hasResources = (
            (selectedNode.value.videos && selectedNode.value.videos.length > 0) ||
            (selectedNode.value.exercises && selectedNode.value.exercises.length > 0) ||
            (selectedNode.value.courseware && selectedNode.value.courseware.length > 0)
          );
          
          if (hasResources) {
            message.success({ content: '已加载知识点详情', key: 'pointDetails' });
          } else {
            message.info({ content: '已加载知识点详情，但暂无关联资源', key: 'pointDetails' });
          }
          
          // 显示抽屉
          drawerVisible.value = true;
        } else {
          console.error('知识点详情加载失败');
          message.error({ content: '知识点详情加载失败', key: 'pointDetails' });
        }
      } catch (error) {
        console.error('知识点详情加载错误:', error);
        message.error({ content: '知识点详情加载错误: ' + (error.message || '未知错误'), key: 'pointDetails' });
      }
    }
  } catch (error) {
    console.error('处理节点点击错误:', error);
    hasError.value = true;
    errorMessage.value = '处理节点点击时发生错误';
    message.error(`加载失败: ${error.message || '未知错误'}`);
  }
};

// 返回上一级
const goBackToParent = () => {
  if (graphStore.currentLevel === 'point' && graphStore.currentChapter) {
    // 从知识点层级返回到小节层级
    graphStore.loadSectionGraph(graphStore.currentChapter.id)
  } else if (graphStore.currentLevel === 'section') {
    // 从小节层级返回到章节层级
    graphStore.loadChapterGraph()
  }
}

// 返回顶层章节
const goToChapterLevel = () => {
  graphStore.loadChapterGraph()
}

// 搜索方法
const onSearch = async (value) => {
  try {
    if (!value || value.trim() === '') {
      message.info('请输入搜索关键词')
      return
    }
    
    message.loading({ content: '正在搜索...', key: 'search' })
    
    const results = await graphStore.searchNodes(value)
    
    if (results && results.length > 0) {
      message.success({ content: `找到 ${results.length} 个相关节点`, key: 'search' })
      searchResults.value = results
      searchResultsVisible.value = true
    } else {
      // 未找到结果
      message.info({ content: '未找到匹配的节点', key: 'search' })
      searchResults.value = []
      searchResultsVisible.value = true
    }
  } catch (error) {
    console.error('搜索错误:', error)
    message.error({ content: '搜索过程中发生错误', key: 'search' })
  }
}

// 搜索结果中选择节点后的导航
const navigateToNode = async (node) => {
  try {
    if (!node || !node.id) {
      console.warn('无效的节点数据:', node)
      return
    }
    
    searchResultsVisible.value = false
    
    if (node.type === 'chapter') {
      // 导航到章节
      await graphStore.loadSectionGraph(node.id)
    } else if (node.type === 'section') {
      // 导航到小节
      await graphStore.loadPointGraph(node.id)
    } else if (node.type === 'point') {
      // 导航到知识点并显示详情
      if (node.sectionId) {
        await graphStore.loadPointGraph(node.sectionId)
      }
      await graphStore.loadPointDetails(node.id)
      selectedNode.value = graphStore.selectedPoint
      drawerVisible.value = true
    }
  } catch (error) {
    console.error('导航到节点错误:', error)
    message.error('导航到节点失败')
  }
}

// 获取节点类型对应的颜色
const getNodeTypeColor = (type) => {
  switch (type) {
    case 'chapter':
      return 'blue'
    case 'section':
      return 'green'
    case 'point':
      return 'orange'
    default:
      return 'default'
  }
}

// 获取与当前节点相关的节点
const relatedNodes = ref([])

// 在handleNodeClick函数中添加相关节点的获取逻辑
const updateRelatedNodes = (node) => {
  if (!node || !graphStore.currentGraphData) return
  
  relatedNodes.value = []
  
  // 查找与当前节点有连接的其他节点
  if (graphStore.currentGraphData.links && graphStore.currentGraphData.nodes) {
    const relatedLinks = graphStore.currentGraphData.links.filter(
      link => link.source === node.id || link.target === node.id
    )
    
    const relatedNodeIds = new Set()
    relatedLinks.forEach(link => {
      if (link.source === node.id) {
        relatedNodeIds.add(link.target)
      } else {
        relatedNodeIds.add(link.source)
      }
    })
    
    relatedNodes.value = graphStore.currentGraphData.nodes.filter(
      n => relatedNodeIds.has(n.id) && n.id !== node.id
    )
  }
}

// 获取章节包含的小节数量
const getChapterSectionCount = (chapterId) => {
  const sections = mockData.value?.sections?.[chapterId] || []
  return sections.length
}

// 获取小节包含的知识点数量
const getSectionPointCount = (sectionId) => {
  const points = mockData.value?.points?.[sectionId] || []
  return points.length
}

// 在watch中添加对selectedNode的监听，更新相关节点
watch(selectedNode, (newNode) => {
  if (newNode) {
    updateRelatedNodes(newNode)
  } else {
    relatedNodes.value = []
  }
})

// 清理选中状态
watch(drawerVisible, (visible) => {
  if (!visible) {
    graphStore.clearSelection()
  }
})

// 添加获取学习状态的方法（模拟，实际应从用户数据中获取）
const getLearningStatus = (pointId) => {
  // 模拟从localStorage获取学习状态，实际项目中这应该从用户数据或后端API获取
  const statusMap = JSON.parse(localStorage.getItem('learningStatus') || '{}')
  return statusMap[pointId] || 'not_started'
}

// 修改updateLearningStatus方法，确保更新后UI状态正确刷新
const updateLearningStatus = (pointId, status) => {
  try {
    console.log('更新学习状态:', pointId, status);
    
    // 更新localStorage中的状态
    const statusMap = JSON.parse(localStorage.getItem('learningStatus') || '{}');
    statusMap[pointId] = status;
    localStorage.setItem('learningStatus', JSON.stringify(statusMap));
    
    // 更新当前选中节点的状态，确保抽屉中显示的状态同步
    if (selectedNode.value && selectedNode.value.id === pointId) {
      // 这里无需直接修改selectedNode，因为getLearningStatus会从localStorage读取最新状态
      // 仅需强制触发UI刷新
      selectedNode.value = { ...selectedNode.value };
    }
    
    // 重新渲染图谱以更新节点状态显示
    updateGraph();
    
    message.success('学习状态已更新');
  } catch (error) {
    console.error('更新学习状态失败:', error);
    message.error('更新学习状态失败');
  }
};

// 预览视频
const previewVideo = (video) => {
  console.log('预览视频:', video);
  
  // 获取正确的URL字段
  const videoUrl = video.play_url || video.url;
  
  if (!videoUrl) {
    message.error('视频URL不存在，无法预览');
    return;
  }
  
  // 确保URL是完整的
  let url = videoUrl;
  if (url && !url.startsWith('http://') && !url.startsWith('https://')) {
    const baseUrl = window.location.origin;
    url = url.startsWith('/') ? `${baseUrl}${url}` : `${baseUrl}/${url}`;
    console.log('完整视频URL:', url);
  }
  
  // 处理封面URL
  let coverUrl = video.cover_url || null;
  if (coverUrl && !coverUrl.startsWith('http://') && !coverUrl.startsWith('https://')) {
    const baseUrl = window.location.origin;
    coverUrl = coverUrl.startsWith('/') ? `${baseUrl}${coverUrl}` : `${baseUrl}/${coverUrl}`;
    console.log('完整封面URL:', coverUrl);
  }
  
  previewModalTitle.value = video.title || '视频预览';
  previewContentType.value = 'video';
  previewContent.value = { ...video, url, cover_url: coverUrl };
  previewModalVisible.value = true;
}

// 预览课件
const previewCourseware = (courseware) => {
  console.log('预览课件:', courseware);
  
  // 获取正确的URL字段
  const coursewareUrl = courseware.courseware_url || courseware.url;
  
  if (!coursewareUrl) {
    message.error('课件URL不存在，无法预览');
    return;
  }
  
  // 确保URL是完整的
  let url = coursewareUrl;
  if (url && !url.startsWith('http://') && !url.startsWith('https://')) {
    const baseUrl = window.location.origin;
    url = url.startsWith('/') ? `${baseUrl}${url}` : `${baseUrl}/${url}`;
    console.log('完整课件URL:', url);
  }
  
  previewModalTitle.value = courseware.title || '课件预览';
  previewContentType.value = url.toLowerCase().endsWith('.pdf') ? 'pdf' : 'other';
  previewContent.value = { ...courseware, url };
  previewModalVisible.value = true;
}

// 预览练习题
const previewExercise = (exercise) => {
  console.log('预览练习题:', exercise);
  
  // 获取正确的URL字段
  const exerciseUrl = exercise.exercise_url || exercise.url;
  
  if (!exerciseUrl) {
    message.error('练习题URL不存在，无法预览');
    return;
  }
  
  // 确保URL是完整的
  let url = exerciseUrl;
  if (url && !url.startsWith('http://') && !url.startsWith('https://')) {
    const baseUrl = window.location.origin;
    url = url.startsWith('/') ? `${baseUrl}${url}` : `${baseUrl}/${url}`;
    console.log('完整练习题URL:', url);
  }
  
  previewModalTitle.value = exercise.title || '练习题预览';
  previewContentType.value = url.toLowerCase().endsWith('.pdf') ? 'pdf' : 'other';
  previewContent.value = { ...exercise, url };
  previewModalVisible.value = true;
}

// 预览模态框
const previewModalVisible = ref(false)
const previewModalTitle = ref('')
const previewContentType = ref('')
const previewContent = ref(null)
const videoPlayer = ref(null)

// 添加预览模态框关闭处理函数
const handlePreviewModalClose = () => {
  // 如果是视频，停止播放
  if (previewContentType.value === 'video' && videoPlayer.value) {
    videoPlayer.value.pause();
    videoPlayer.value.currentTime = 0;
  }
  
  // 清空预览内容
  previewContent.value = null;
  previewContentType.value = '';
}
</script>

<style scoped>
.knowledge-graph {
}
.difficulty-assessment {
  padding: 8px;
}

.difficulty-score {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.score-label {
  margin-right: 16px;
  font-weight: 500;
  color: rgba(0, 0, 0, 0.85);
}

.score-value {
  margin-left: 16px;
  color: rgba(0, 0, 0, 0.45);
}

.difficulty-progress {
  margin-bottom: 16px;
}

.difficulty-stats {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-top: 1px solid #f0f0f0;
}

.stat-item {
  flex: 1;
  text-align: center;
}

.stat-label {
  color: rgba(0, 0, 0, 0.45);
  font-size: 12px;
  margin-bottom: 4px;
}

.stat-value {
  color: rgba(0, 0, 0, 0.85);
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.knowledge-graph {
  width: 100%;
  height: calc(100vh - 250px);
  position: relative;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #f9f9f9 0%, #f0f2f5 100%);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.graph-toolbar {
  padding: 16px;
  background: linear-gradient(to right, rgba(255, 255, 255, 0.95), rgba(255, 255, 255, 0.85));
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  display: flex;
  justify-content: center;
  align-items: center;
}

.graph-toolbar :deep(.ant-input-search) {
  width: 280px;
  
  .ant-input {
    height: 40px;
    padding: 8px 16px;
    font-size: 14px;
    border-radius: 8px 0 0 8px;
    border: 1px solid rgba(0, 0, 0, 0.1);
    border-right: none;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(8px);
    transition: all 0.3s ease;
    
    &:hover {
      border-color: #4CAF50;
      box-shadow: 0 2px 8px rgba(33, 150, 243, 0.15);
    }
    
    &:focus {
      border-color: #2196F3;
      box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2);
      outline: none;
    }
  }
  
  .ant-input-search-button {
    height: 40px;
    width: 48px;
    border-radius: 0 8px 8px 0;
    border: none;
    background: linear-gradient(120deg, #2196F3, #4CAF50);
    transition: all 0.3s ease;
    margin-left: -1px;
    box-shadow: 0 2px 8px rgba(33, 150, 243, 0.2);
    
    &:hover {
      background: linear-gradient(135deg, #40a9ff 0%, #1890ff 100%);
      box-shadow: 0 2px 8px rgba(24, 144, 255, 0.25);
    }
    
    &:active {
      background: linear-gradient(135deg, #096dd9 0%, #0050b3 100%);
    }
  }
}

.graph-toolbar :deep(.ant-btn-group) {
  .ant-btn {
    height: 40px;
    width: 40px;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #fff;
    border: 1px solid rgba(0, 0, 0, 0.1);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;
    color: #5c6b7c;
    
    &:hover {
      color: #2196F3;
      border-color: #2196F3;
      background: rgba(255, 255, 255, 0.95);
      transform: translateY(-1px);
      box-shadow: 0 2px 8px rgba(33, 150, 243, 0.15);
      
      .anticon {
        transform: scale(1.1);
      }
    }
    
    &:active {
      transform: translateY(1px);
      box-shadow: 0 1px 2px rgba(33, 150, 243, 0.25);
    }
    
    .anticon {
      font-size: 16px;
      transition: transform 0.3s ease;
    }
    
    &:first-child {
      border-top-left-radius: 8px;
      border-bottom-left-radius: 8px;
    }
    
    &:last-child {
      border-top-right-radius: 8px;
      border-bottom-right-radius: 8px;
    }
  }
}

.graph-toolbar :deep(.ant-space) {
  gap: 16px !important;
}
.graph-breadcrumb {
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(8px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  font-size: 14px;
}

.graph-breadcrumb a {
  color: #2196F3;
  transition: all 0.3s ease;
  font-weight: 500;
}

.graph-breadcrumb a:hover {
  color: #1976D2;
  text-decoration: underline;
}

.graph-container {
  flex: 1;
  overflow: hidden;
  background: linear-gradient(135deg, #f5f7fa 0%, #eef2f7 100%);
  position: relative;
  border-radius: 0 0 16px 16px;
  box-shadow: none;
  transition: all 0.3s ease;
}

/* 节点样式 */
:deep(.node) {
  cursor: pointer;
  transition: all 0.3s ease;
}

:deep(.node:hover) circle {
  stroke: #1890ff;
  stroke-width: 4px;
  filter: drop-shadow(0 0 8px rgba(24, 144, 255, 0.5));
}

:deep(.link) {
  pointer-events: none;
  transition: stroke-width 0.2s ease;
  opacity: 0.7;
}

:deep(.link:hover) {
  stroke-width: 3px;
  opacity: 1;
}

/* 空数据提示 */
.empty-data {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  width: 100%;
  animation: fadeIn 0.5s ease-out;
  max-width: 300px;
  background: rgba(255, 255, 255, 0.9);
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

/* 错误提示 */
.error-message {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  padding: 24px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  text-align: center;
  max-width: 80%;
  z-index: 100;
  backdrop-filter: blur(10px);
}

.error-message .title {
  font-size: 18px;
  color: #ff4d4f;
  margin-bottom: 12px;
  font-weight: bold;
}

.error-message .content {
  color: #666;
  margin-bottom: 16px;
  line-height: 1.6;
}

.legend-panel {
  position: absolute;
  right: 15px;
  top: 80px;
  background: rgba(255, 255, 255, 0.95);
  padding: 15px;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  z-index: 10;
  max-width: 220px;
  backdrop-filter: blur(12px);
  transition: all 0.3s ease;
}

.legend-panel:hover {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.legend-panel h3 {
  margin: 5px 0 10px 0;
  font-size: 16px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  padding-bottom: 8px;
  color: #2196F3;
  font-weight: 600;
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.legend-items {
  margin-bottom: 15px;
}

.legend-item {
  display: flex;
  align-items: center;
  margin: 8px 0;
  flex-wrap: wrap;
  transition: all 0.3s ease;
  padding: 4px 0;
  border-radius: 4px;
}

.legend-item:hover {
  background-color: rgba(33, 150, 243, 0.05);
  transform: translateX(2px);
}

.legend-color {
  display: inline-block;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  margin-right: 10px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.legend-label {
  font-weight: 600;
  margin-right: 8px;
  color: #262626;
}

.legend-text {
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 4px;
  margin-left: 26px;
  width: 100%;
}

/* 知识点详情样式 */
.knowledge-detail {
  padding-bottom: 20px;
}

.knowledge-title {
  font-size: 22px;
  font-weight: bold;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.detail-card {
  margin-bottom: 20px;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.95);
}

.detail-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.section-title {
  font-size: 16px;
  margin: 0;
  color: #2196F3;
  display: flex;
  align-items: center;
  font-weight: 600;
}

.section-title .anticon {
  margin-right: 8px;
  color: #2196F3;
}

.resource-count {
  font-size: 14px;
  color: #8c8c8c;
  margin-left: 8px;
  font-weight: normal;
}

.resource-list-item {
  border-bottom: 1px solid rgba(0, 0, 0, 0.04);
  padding: 10px 0;
  transition: all 0.3s ease;
}

.resource-list-item:hover {
  background-color: rgba(33, 150, 243, 0.05);
}

.resource-list-item:last-child {
  border-bottom: none;
}

.resource-title {
  font-weight: 500;
  margin-bottom: 4px;
}

.resource-title a {
  color: #2196F3;
  transition: all 0.3s ease;
}

.resource-title a:hover {
  color: #1976D2;
  text-decoration: underline;
}

.resource-description {
  font-size: 13px;
  color: #595959;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  line-height: 1.6;
}

.description-content {
  white-space: pre-line;
  padding: 12px;
  background-color: rgba(0, 0, 0, 0.02);
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.8;
  color: #595959;
  transition: all 0.3s ease;
}

.no-resources {
  margin-top: 24px;
  padding: 20px;
  background-color: rgba(0, 0, 0, 0.02);
  border-radius: 12px;
  text-align: center;
  transition: all 0.3s ease;
}

/* 适配暗色主题 */
:deep([data-theme='dark']) {
  .knowledge-graph {
    background: linear-gradient(135deg, #1a1f2c 0%, #2d3748 100%);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  }
  
  .graph-toolbar,
  .graph-breadcrumb {
    background: rgba(26, 32, 44, 0.8);
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .graph-container {
    background: linear-gradient(135deg, #1a1f2c 0%, #2d3748 100%);
  }
  
  .legend-panel {
    background: rgba(26, 32, 44, 0.9);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  }
  
  .legend-panel h3 {
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    color: #e2e8f0;
  }
  
  .legend-label {
    color: #e2e8f0;
  }
  
  .legend-text {
    color: #a0aec0;
  }
  
  .empty-data,
  .error-message {
    background: rgba(26, 32, 44, 0.9);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  }
  
  .error-message .content {
    color: #a0aec0;
  }
  
  .knowledge-title {
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    background: linear-gradient(90deg, #60a5fa, #34d399);
    -webkit-background-clip: text;
    background-clip: text;
  }
  
  .detail-card {
    background: rgba(26, 32, 44, 0.8);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .section-title {
    color: #e2e8f0;
  }
  
  .resource-list-item {
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .resource-list-item:hover {
    background-color: rgba(96, 165, 250, 0.1);
  }
  
  .resource-title a {
    color: #60a5fa;
  }
  
  .resource-title a:hover {
    color: #93c5fd;
  }
  
  .resource-description {
    color: #a0aec0;
  }
  
  .description-content {
    background-color: rgba(255, 255, 255, 0.05);
    color: #a0aec0;
  }
  
  .no-resources {
    background-color: rgba(255, 255, 255, 0.05);
  }
}
</style>