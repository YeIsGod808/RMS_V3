import { useUserStore } from '@/store/modules/user'

export function setupPermission(app) {
  app.directive('permission', {
    mounted(el, binding) {
      const userStore = useUserStore()
      const requiredRole = binding.value
      
      if (requiredRole === 'admin' && !userStore.isAdmin) {
        el.style.display = 'none'
        
        el.dataset.permissionHidden = 'true'
      }
    }
  })
} 