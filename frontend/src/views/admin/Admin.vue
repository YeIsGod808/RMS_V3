<template>
  <div class="admin-page">
    <div class="page-header">
      <div class="header-content">
        <h1 class="title">管理后台</h1>
        <div class="subtitle">系统配置与数据管理中心</div>
      </div>
    </div>
    <div class="card-container">
      <Loading :loading="adminStore.loading" />
      <a-tabs v-model:activeKey="activeKey" class="admin-tabs">
        <!-- 用户管理 -->
        <a-tab-pane key="users" tab="用户管理" v-permission="'admin'">
          <UserManagement />
        </a-tab-pane>
        
        <!-- 用户组管理 -->
        <a-tab-pane key="userGroups" tab="用户组管理" v-permission="'admin'">
          <UserGroupManagement />
        </a-tab-pane>
        
        <!-- 资源管理 -->
        <a-tab-pane key="resources" tab="资源管理">
          <ResourceManagement />
        </a-tab-pane>
        
        <!-- 知识图谱管理 -->
        <a-tab-pane key="knowledge" tab="知识图谱管理">
          <a-tabs>
            <a-tab-pane key="nodes" tab="节点管理">
              <KnowledgeNodeManagement />
            </a-tab-pane>
            <a-tab-pane key="relations" tab="关系管理">
              <KnowledgeRelationManagement />
            </a-tab-pane>
            <a-tab-pane key="import" tab="一键导入知识图谱">
              <KnowledgeImport />
            </a-tab-pane>
          </a-tabs>
        </a-tab-pane>
      </a-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAdminStore } from '@/store/modules/admin'
import { useKnowledgeAdminStore } from '@/store/modules/knowledgeAdmin'
import Loading from '@/components/Loading.vue'
import UserManagement from './UserManagement.vue'
import ResourceManagement from './ResourceManagement.vue'
import KnowledgeNodeManagement from './KnowledgeNodeManagement.vue'
import KnowledgeRelationManagement from './KnowledgeRelationManagement.vue'
import UserGroupManagement from './UserGroupManagement.vue'
import KnowledgeImport from './KnowledgeImport.vue'

const adminStore = useAdminStore()
const knowledgeAdminStore = useKnowledgeAdminStore()
const activeKey = ref('users')

onMounted(async () => {
  //await adminStore.getUsers()
})
</script>

<style scoped>
.admin-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f0f7ff 0%, #e0eafc 100%);
  padding: 24px;
  transition: background 0.3s ease;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 24px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(12px);
  transition: all 0.3s ease;
  animation: slideDown 0.5s ease-out;
}

@keyframes slideDown {
  0% {
    opacity: 0;
    transform: translateY(-20px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

.header-content .title {
  font-size: 32px;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  -webkit-background-clip: text;
  /* 添加标准属性 */
  background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 0.5px;
}

.header-content .subtitle {
  color: #5c6b7c;
  margin-top: 8px;
  font-size: 16px;
  letter-spacing: 0.3px;
}

.card-container {
  position: relative;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(12px);
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.2);
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  0% {
    opacity: 0;
    transform: translateY(10px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

.admin-tabs {
  margin-top: 10px;
}

:deep(.ant-tabs-nav) {
  margin-bottom: 20px;
  background: rgba(255, 255, 255, 0.8);
  padding: 8px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

:deep(.ant-tabs-tab) {
  padding: 12px 20px;
  transition: all 0.3s ease;
  border-radius: 8px;
  margin-right: 8px;
}

:deep(.ant-tabs-tab-active) {
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  border-radius: 8px;
}

:deep(.ant-tabs-tab-active .ant-tabs-tab-btn) {
  color: white !important;
  font-weight: 500;
}

:deep(.ant-tabs-ink-bar) {
  display: none;
}

:deep(.ant-tabs-content) {
  background: rgba(255, 255, 255, 0.7);
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);
}
</style>