<template>
  <a-layout-header class="header">
    <div class="logo">数据库知识图谱</div>
    <a-menu
      v-model:selectedKeys="selectedKeys"
      theme="dark"
      mode="horizontal"
      :style="{ lineHeight: '64px' }"
    >
      <a-menu-item key="graph">
        <router-link to="/graph">知识图谱</router-link>
      </a-menu-item>
      <a-menu-item key="admin" v-if="userStore.isAdmin">
        <router-link to="/admin">管理后台</router-link>
      </a-menu-item>
      <div class="user-menu">
        <a-dropdown>
          <a class="ant-dropdown-link" @click.prevent>
            {{ userStore.username }}
            <down-outlined />
          </a>
          <template #overlay>
            <a-menu>
              <a-menu-item key="logout" @click="handleLogout">
                退出登录
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </div>
    </a-menu>
  </a-layout-header>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/modules/user'
import { DownOutlined } from '@ant-design/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const selectedKeys = ref(['graph'])

const handleLogout = async () => {
  await userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.header {
  display: flex;
  align-items: center;
}

.logo {
  color: white;
  font-size: 18px;
  margin-right: 30px;
}

.user-menu {
  float: right;
  margin-left: auto;
  color: white;
}
</style> 