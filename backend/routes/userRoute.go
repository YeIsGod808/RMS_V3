package routes

import (
	"github.com/RMS_V3/internal/user"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	// 创建用户相关API的路由组
	userGroup := r.Group("/user")
	{
		// 校验用户token
		userGroup.GET("/checktoken", user.CheckToken)
		// 用户登录接口，生成并返回token
		userGroup.POST("/login", user.GenerateToken)
		// 修改用户密码
		userGroup.POST("/change-password", user.ChangePsw)
		// 添加单个用户
		userGroup.POST("/add-user", user.AddUser)
		// 批量导入用户
		userGroup.POST("/batch-import", user.AddUserBatch)
		// 用户注册接口（可根据需求开启或禁用）
		userGroup.POST("/register", user.Register)
		// 获取用户列表
		userGroup.GET("/list", user.ListUsers)
		// 更新用户信息
		userGroup.POST("/update", user.UpdateUser)
		// 删除用户
		userGroup.POST("/delete", user.DeleteUser)
	}

	// 创建用户组相关API的路由组
	user_groupGroup := r.Group("/user-group")
	{
		// 获取用户组下的用户信息
		user_groupGroup.GET("/get-user", user.GetGroupUser)
		// 获取所有用户组列表
		user_groupGroup.GET("/get-groups", user.GetGroupList)
		// 创建新的用户组
		user_groupGroup.POST("/create", user.CreateGroup)
		// 向用户组添加用户
		user_groupGroup.POST("/add-user", user.AddGroupUser)
		// 从用户组中删除用户
		user_groupGroup.POST("/delete-user", user.DeleteGroupUser)
		// 删除用户组
		user_groupGroup.POST("/delete-group", user.DeleteGroup)
		// 获取所有用户组信息
		user_groupGroup.GET("/all-groups", user.GetAllGroup)
		// 编辑用户组名称
		user_groupGroup.POST("/edit-name", user.EditGroupName)
		// 检查是否为用户组所有者（可根据需求开启或禁用）
		//user_groupGroup.GET("/if-is-owner", IsGroupOwner)
	}
}
