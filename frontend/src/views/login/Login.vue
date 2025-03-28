<template>
  <div class="login-container">
    <a-card class="login-card" title="数据库教学资源管理系统">
      <a-form
        :model="formState"
        @finish="handleFinish"
        class="form-container"
      >
        <a-form-item
          name="username"
          :rules="[{ required: true, message: '请输入用户名' }]"
        >
          <a-input v-model:value="formState.username" placeholder="用户名" autocomplete="off">
            <template #prefix>
              <user-outlined />
            </template>
          </a-input>
        </a-form-item>
        
        <a-form-item
          name="password"
          :rules="[{ required: true, message: '请输入密码' }]"
        >
          <a-input-password v-model:value="formState.password" placeholder="密码" autocomplete="off">
            <template #prefix>
              <lock-outlined />
            </template>
          </a-input-password>
        </a-form-item>
        
        <a-form-item>
          <a-button type="primary" html-type="submit" block :loading="loading">
            登录
          </a-button>
          <div class="register-link">
            没有账号？<router-link to="/register">立即注册</router-link>
          </div>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { useUserStore } from '@/store/modules/user'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)

const formState = reactive({
  username: '',
  password: ''
})

const handleFinish = async (values) => {
  try {
    loading.value = true
    await userStore.login(values.username, values.password)
    message.success('登录成功')
    
    // 如果有重定向参数，则重定向到该路径，否则跳转到首页
    const redirect = router.currentRoute.value.query.redirect
    router.push(redirect || '/')
  } catch (error) {
    // 显示更具体的错误信息
    if (error.message === 'wrong password') {
      message.error('密码错误')
    } else if (error.message === 'user not found') {
      message.error('用户不存在')
    } else {
      message.error(error.message || '登录失败')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(-45deg, #ee7752, #e73c7e, #23a6d5, #23d5ab);
  background-size: 400% 400%;
  animation: gradient 15s ease infinite;
  overflow: hidden;
}

@keyframes gradient {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.login-card {
  width: 368px;
  margin-top: -50px;
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.9) !important;
  border-radius: 16px !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  animation: cardFloat 0.8s ease-out;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.login-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 36px rgba(0, 0, 0, 0.15);
}

@keyframes cardFloat {
  0% {
    opacity: 0;
    transform: translateY(20px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

:deep(.ant-card-head) {
  text-align: center;
  border-bottom: none;
  padding-bottom: 0;
}

:deep(.ant-card-head-title) {
  font-size: 28px;
  font-weight: 600;
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  -webkit-background-clip: text;
  /* 添加标准属性 */
  background-clip: text;
  -webkit-text-fill-color: transparent;
  padding: 16px 0;
}

:deep(.ant-card-body) {
  padding: 32px 24px;
}

:deep(.ant-input-affix-wrapper) {
  border-radius: 12px;
  padding: 8px 12px;
  border: 2px solid transparent;
  background-color: rgba(255, 255, 255, 0.8);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

:deep(.ant-input-affix-wrapper:hover),
:deep(.ant-input-affix-wrapper-focused) {
  border-color: #4CAF50;
  background-color: rgba(255, 255, 255, 0.95);
  box-shadow: 0 4px 12px rgba(76, 175, 80, 0.15);
  transform: translateY(-1px);
}

:deep(.ant-input) {
  background: transparent;
  font-size: 15px;
}

:deep(.ant-input-prefix) {
  margin-right: 12px;
  color: #4CAF50;
  opacity: 0.8;
  transition: opacity 0.3s ease;
}

:deep(.ant-input-affix-wrapper:hover .ant-input-prefix),
:deep(.ant-input-affix-wrapper-focused .ant-input-prefix) {
  opacity: 1;
}

:deep(.ant-btn-primary) {
  height: 40px;
  border-radius: 8px;
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  border: none;
  transition: all 0.3s ease;
}

:deep(.ant-btn-primary:hover) {
  background: linear-gradient(120deg, #1976D2, #388E3C);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.register-link {
  margin-top: 16px;
  text-align: center;
  opacity: 0;
  animation: fadeIn 0.5s ease forwards 0.8s;
}

@keyframes fadeIn {
  to {
    opacity: 1;
  }
}

:deep(.ant-form-item) {
  margin-bottom: 24px;
  opacity: 0;
  animation: slideIn 0.5s ease forwards;
}

:deep(.ant-form-item:nth-child(1)) { animation-delay: 0.2s; }
:deep(.ant-form-item:nth-child(2)) { animation-delay: 0.4s; }
:deep(.ant-form-item:nth-child(3)) { animation-delay: 0.6s; }

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}
</style>