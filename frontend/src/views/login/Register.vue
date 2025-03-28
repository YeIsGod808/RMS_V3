<template>
  <div class="register-container">
    <a-card class="register-card" title="用户注册">
      <a-form
        :model="formState"
        :rules="rules"
        @finish="handleFinish"
        class="form-container"
      >
        <a-form-item name="username">
          <a-input v-model:value="formState.username" placeholder="用户名">
            <template #prefix>
              <user-outlined />
            </template>
          </a-input>
        </a-form-item>
        
        <a-form-item name="nickname">
          <a-input v-model:value="formState.nickname" placeholder="昵称">
            <template #prefix>
              <smile-outlined />
            </template>
          </a-input>
        </a-form-item>
        
        <a-form-item name="password">
          <a-input-password v-model:value="formState.password" placeholder="密码">
            <template #prefix>
              <lock-outlined />
            </template>
          </a-input-password>
        </a-form-item>
        
        <a-form-item name="confirmPassword">
          <a-input-password v-model:value="formState.confirmPassword" placeholder="确认密码">
            <template #prefix>
              <safety-outlined />
            </template>
          </a-input-password>
        </a-form-item>
        
        <a-form-item>
          <a-button type="primary" html-type="submit" block :loading="loading">
            注册
          </a-button>
          <div class="login-link">
            已有账号？<router-link to="/login">立即登录</router-link>
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
import { register } from '@/api/user'
import { UserOutlined, LockOutlined, SafetyOutlined, SmileOutlined } from '@ant-design/icons-vue'
import { errcode } from '@/utils/errcode'

const router = useRouter()
const loading = ref(false)

const formState = reactive({
  username: '',
  nickname: '',
  password: '',
  confirmPassword: ''
})

// 表单验证规则
const rules = {
  username: [
    { required: true, message: '请输入用户名' },
    { pattern: /^([a-z]|[A-Z]|[0-9]|_){1,32}$/, message: '用户名只能包含字母、数字和下划线，长度1-32位' }
  ],
  nickname: [
    { required: true, message: '请输入昵称' },
    { max: 128, message: '昵称长度不能超过128个字符' }
  ],
  password: [
    { required: true, message: '请输入密码' },
    { pattern: /^([a-z]|[A-Z]|[0-9]){8,16}$/, message: '密码只能包含字母和数字，长度8-16位' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码' },
    { validator: validateConfirmPassword }
  ]
}

// 验证两次密码是否一致
function validateConfirmPassword(rule, value) {
  if (value !== formState.password) {
    return Promise.reject('两次输入的密码不一致')
  }
  return Promise.resolve()
}

const handleFinish = async (values) => {
  try {
    loading.value = true
    const response = await register(values.username, values.password, values.nickname)
    
    if (response.ret === '0') {
      message.success('注册成功，请登录')
      router.push('/login')
    } else {
      // 处理不同类型的错误
      switch (response.ret) {
        case '1':
          if (response.msg === 'invalid register info') {
            message.error('注册信息无效，请检查用户名和密码格式')
          } else {
            message.error(response.msg || '注册失败')
          }
          break
        case errcode.DB_DUP_ERR:
          message.error('用户名或昵称已被使用，请更换后重试')
          break
        case errcode.DB_CONN_ERR:
          message.error('服务器错误，请稍后重试')
          break
        default:
          message.error(response.msg || '注册失败')
      }
    }
  } catch (error) {
    console.error('注册错误:', error)
    // 处理网络错误或其他未预期的错误
    if (error.response) {
      // 服务器返回了错误状态码
      switch (error.response.status) {
        case 400:
          if (error.response.data.msg === 'id or nickname is used') {
            message.error('用户名或昵称已被使用，请更换后重试')
          } else if (error.response.data.msg === 'invalid register info') {
            message.error('注册信息无效，请检查输入格式')
          } else {
            message.error(error.response.data.msg || '注册参数无效')
          }
          break
        case 500:
          message.error('服务器错误，请稍后重试')
          break
        default:
          message.error('注册失败，请稍后重试')
      }
    } else if (error.request) {
      // 请求发出但没有收到响应
      message.error('网络连接失败，请检查网络后重试')
    } else {
      // 请求设置时发生错误
      message.error('注册失败：' + error.message)
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
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

.register-card {
  width: 368px;
  margin-top: -50px;
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.9) !important;
  border-radius: 16px !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  animation: cardFloat 0.8s ease-out;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.register-card:hover {
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

.login-link {
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
:deep(.ant-form-item:nth-child(4)) { animation-delay: 0.8s; }

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