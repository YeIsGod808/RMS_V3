package user

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/RMS_V3/log/logger"
	"github.com/RMS_V3/pkg/commonlib"
	"github.com/RMS_V3/pkg/errcode"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func BatchAddUserToGroup(c *gin.Context) {
	// 检查用户权限
	u, ok := CheckUserPermission(c.Query("token"), Admin, c)
	if !ok {
		return
	}

	// 获取 group_id 并验证
	groupID, err := strconv.Atoi(c.Query("group_id"))
	if err != nil || groupID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"ret": http.StatusBadRequest,
			"msg": "group_id为空或格式错误",
		})
		return
	}

	// 解析上传的文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		logger.ERROR_LOG(err.Error(), c)
		c.JSON(http.StatusBadRequest, gin.H{
			"ret": errcode.WRONG_PARAM,
			"msg": "上传失败，请检查文件格式",
		})
		return
	}

	// 打开文件并解析 CSV
	file, err := fileHeader.Open()
	if err != nil {
		logger.ERROR_LOG(err.Error(), c)
		c.JSON(http.StatusBadRequest, gin.H{
			"ret": errcode.WRONG_PARAM,
			"msg": "文件解析失败，请检查文件格式",
		})
		return
	}
	defer file.Close()

	// 初始化 CSV 读取器
	reader := csv.NewReader(bufio.NewReader(file))

	// 默认密码和允许的用户类型
	defaultPassword := md5Gen("12345678")
	validUserTypes := map[string]bool{"student": true}
	if u.UserType == Admin {
		validUserTypes["teacher"] = true
		validUserTypes["admin"] = true
	}

	// 处理 CSV 文件内容
	users, err := processCSV(reader, validUserTypes, u.UserType)
	if err != nil {
		logger.ERROR_LOG(err.Error(), c)
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.WRONG_PARAM,
			"msg": err.Error(),
		})
		return
	}

	// 插入用户数据到数据库
	err = insertUsersIntoDB(users, defaultPassword)
	if err != nil {
		logger.ERROR_LOG(err.Error(), c)
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.DBERR_BASE,
			"msg": "插入账号时出现错误，请检查文件内容",
		})
		return
	}

	// 将用户添加到组
	if groupID != 0 {
		valid, err := AddGroupUserDB(u, groupID, users)
		if !valid {
			c.JSON(http.StatusForbidden, gin.H{
				"ret": errcode.AUTH_ERR,
				"msg": fmt.Sprintf("导入失败，用户不允许操作组:%d", groupID),
			})
			return
		}
		if err != nil {
			if commonlib.IsDbDupErr(err) {
				c.JSON(http.StatusBadGateway, gin.H{
					"ret": errcode.DBERR_BASE,
					"msg": "导入失败，表单中存在组内已有账户",
				})
				return
			}
			logger.ERROR_LOG(err.Error(), c)
			c.JSON(http.StatusInternalServerError, gin.H{
				"ret": errcode.DBERR_BASE,
				"msg": "服务器内部错误",
			})
			return
		}
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "ok",
	})
}
func processCSV(reader *csv.Reader, validUserTypes map[string]bool, userType UserTypes) ([]string, error) {
	var users []string

	for i := 0; ; i++ {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("第%d行读取失败，请检查文件内容", i+1)
		}

		// 跳过表头
		if i == 0 {
			continue
		}

		// 验证字段数量
		if (userType == Admin && len(line) != 3) || (userType == Teacher && len(line) != 2) {
			return nil, fmt.Errorf("第%d行有缺失或多余字段，请检查文件内容", i+1)
		}

		// 教师用户统一只能添加学生类型
		if userType == Teacher {
			line = append(line, "student")
		}

		// 验证字段内容
		for j := 0; j < 3; j++ {
			if !utf8.ValidString(line[j]) {
				return nil, fmt.Errorf("第%d行存在非UTF-8编码字符，请将文件编码转为UTF-8", i+1)
			}
		}
		if !validUserTypes[line[2]] {
			return nil, fmt.Errorf("第%d行存在不允许当前用户添加的账号类型:%s", i+1, line[2])
		}
		if !useridRegexp.MatchString(line[0]) {
			return nil, fmt.Errorf("第%d行用户账号非法，只允许使用32位内数字，字母，下划线", i+1)
		}

		users = append(users, line[0])
	}

	return users, nil
}
func insertUsersIntoDB(users []string, defaultPassword string) error {
	sqlTemplate := "INSERT INTO t_account (user_id, nickname, password, user_type) VALUES (?, ?, ?, ?)"

	for _, user := range users {
		tx, err := commonlib.DB_user.Begin()
		if err != nil {
			return fmt.Errorf("无法开启事务: %v", err)
		}

		_, err = tx.Exec(sqlTemplate, user[0], user[1], defaultPassword, user[2])
		if err != nil {
			tx.Rollback()
			mysqlErr, ok := err.(*mysql.MySQLError)
			if ok && mysqlErr.Number == 1062 {
				continue // 忽略重复记录
			}
			return fmt.Errorf("插入用户数据失败: %v", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("提交事务失败: %v", err)
		}
	}

	return nil
}
