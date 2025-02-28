package user

import (
	"bufio"
	"crypto/md5"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/RMS_V3/commonlib"
	"github.com/RMS_V3/logger"
	"github.com/RMS_V3/pkg/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var useridRegexp *regexp.Regexp
var passwordRegexp *regexp.Regexp

func init() {
	useridRegexp = regexp.MustCompile("^([a-z]|[A-Z]|[0-9]|_){1,32}$")
	passwordRegexp = regexp.MustCompile("^([a-z]|[A-Z]|[0-9]){8,16}$")
}

func md5Gen(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

type UserTypes string

const (
	Admin   UserTypes = "admin"
	Teacher UserTypes = "teacher"
	Student UserTypes = "student"
)

var typeVal = map[UserTypes]int{Admin: 16, Teacher: 8, Student: 4}

type User struct {
	Id       string    `json:"id"`
	Password string    `json:"password"`
	Nickname string    `json:"nickname"`
	UserType UserTypes `json:"usertype"`
}

func GetUserInfo(token string) (userInfo *User, err error) {
	userInfo = new(User)
	err = nil
	return
}

// GenerateToken godoc
// @Summary 生成用户登录token
// @Description 根据用户ID和密码生成登录token
// @Tags 用户管理
// @Accept  json
// @Produce  json
// @Param id formData string true "用户ID"
// @Param password formData string true "用户密码"
// @Success 200 {object} map[string]interface{} "成功响应"
// @Failure 400 {object} map[string]interface{} "无效参数"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/generate_token [post]
func GenerateToken(c *gin.Context) {
	if !useridRegexp.MatchString(c.PostForm("id")) ||
		!passwordRegexp.MatchString(c.PostForm("password")) {
		logger.DEBUG_LOG("invalid id or psw", c)
		c.JSON(http.StatusBadRequest, gin.H{
			"ret": errcode.MISSING_PARAM,
			"msg": "invalid id or password",
		})
		return
	}

	var userInfo User
	err := commonlib.DB_user.QueryRow("SELECT user_id,nickname,password,user_type FROM t_account WHERE user_id=?", c.PostForm("id")).
		Scan(&userInfo.Id, &userInfo.Nickname, &userInfo.Password, &userInfo.UserType)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.DEBUG_LOG(err.Error(), c)
			c.JSON(http.StatusOK, gin.H{
				"ret": errcode.WRONG_PARAM,
				"msg": "user not found",
			})
			return
		}
		logger.ERROR_LOG(err.Error(), c)
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.DB_CONN_ERR,
			"msg": "Internal Server Error",
		})
		return
	}

	if userInfo.Password == md5Gen(c.PostForm("password")) {
		token, err := jwtGenerateToken(&userInfo, time.Hour*10)
		if err != nil {
			logger.ERROR_LOG(err.Error(), c)
			c.JSON(http.StatusInternalServerError, gin.H{
				"ret": errcode.JWT_ERR,
				"msg": "Internal Server Error",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"ret": "0",
			"msg": "ok",
			"user": map[string]string{
				"id":       userInfo.Id,
				"nickname": userInfo.Nickname,
				"usertype": string(userInfo.UserType),
				"token":    token,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"ret": errcode.WRONG_PASSWORD,
			"msg": "wrong password",
		})
	}
}

func CheckToken(c *gin.Context) {
	userInfo, err := JwtParseToken(c.Query("token"))
	if err != nil {
		logger.DEBUG_LOG(err.Error(), c)
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.JWT_ERR,
			"msg": "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": "0",
		"msg": "ok",
		"user": map[string]string{
			"id":       userInfo.Id,
			"nickname": userInfo.Nickname,
			"usertype": string(userInfo.UserType),
		},
	})
}

// 开放注册接口只能注册学生账号
// 管理员账号通过sql创建
// 教师账号由管理员添加

// Register godoc
// @Summary 注册新用户（仅限学生）
// @Description 开放注册接口只能注册学生账号
// @Tags 用户管理
// @Accept  json
// @Produce  json
// @Param id formData string true "用户ID"
// @Param password formData string true "用户密码"
// @Param nickname formData string true "用户昵称"
// @Success 200 {object} map[string]interface{} "成功响应"
// @Failure 400 {object} map[string]interface{} "无效输入信息"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/register [post]
func Register(c *gin.Context) {
	userInfo := &User{
		Id:       c.PostForm("id"),
		Password: c.PostForm("password"),
		Nickname: c.PostForm("nickname"),
		UserType: Student,
	}
	if !useridRegexp.MatchString(userInfo.Id) ||
		!passwordRegexp.MatchString(userInfo.Password) ||
		len(userInfo.Nickname) > 128 {
		c.JSON(http.StatusBadRequest, gin.H{
			"ret": "1",
			"msg": "invalid register info",
		})
		return
	}
	userInfo.Password = md5Gen(userInfo.Password)

	_, err := commonlib.DB_user.Exec("INSERT INTO t_account (user_id,nickname,password,user_type) VALUES (?,?,?,?)",
		userInfo.Id,
		userInfo.Nickname,
		userInfo.Password,
		userInfo.UserType)

	if err != nil {
		if len(err.Error()) >= 10 && err.Error()[0:10] == "Error 1062" {
			c.JSON(http.StatusBadRequest, gin.H{
				"ret": errcode.DB_DUP_ERR,
				"msg": "id or nickname is used",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.DB_CONN_ERR,
			"msg": "Internal Server Error",
		})
		logger.ERROR_LOG(err.Error(), c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": "0",
		"msg": "succ",
	})
}

// ChangePsw godoc
// @Summary 修改用户密码
// @Description 根据token验证身份后修改用户密码
// @Tags 用户管理
// @Accept  json
// @Produce  json
// @Param token query string true "用户token"
// @Param new_password query string true "新密码"
// @Success 200 {object} map[string]interface{} "成功响应"
// @Failure 400 {object} map[string]interface{} "无效参数或权限不足"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/change_psw [get]
func ChangePsw(c *gin.Context) {
	u, ok := CheckUserPermission(c.Query("token"), Student, c)
	if !ok {
		return
	}
	u.Password = c.Query("new_password")
	if !passwordRegexp.MatchString(u.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"ret": errcode.WRONG_PARAM,
			"msg": "invalid new password",
		})
		return
	}
	_, err := commonlib.DB_user.Exec("UPDATE t_account SET password = ? WHERE user_id = ?", md5Gen(u.Password), u.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.DB_CONN_ERR,
			"msg": "Internal Server Error",
		})
		logger.ERROR_LOG(err.Error(), c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": "0",
		"msg": "succ",
	})
}

// AddUser godoc
// @Summary 批量添加用户
// @Description 根据上传的CSV文件批量添加用户，需要管理员或教师权限
// @Tags 用户管理
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "包含用户信息的CSV文件"
// @Param token query string true "用户token"
// @Param group_id query int false "组ID"
// @Success 200 {object} map[string]interface{} "成功响应"
// @Failure 400 {object} map[string]interface{} "无效参数或文件格式错误"
// @Failure 403 {object} map[string]interface{} "权限不足"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/add_user [post]
func AddUser(c *gin.Context) {
	u, ok := CheckUserPermission(c.Query("token"), Teacher, c)
	if !ok {
		return
	}
	postFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ret": errcode.WRONG_PARAM,
			"msg": "上传失败，请检查文件格式",
		})
		return
	}
	file, err := postFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ret": errcode.WRONG_PARAM,
			"msg": "文件解析失败，请检查文件格式",
		})
		return
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	sqlTeml := "INSERT into t_account (user_id,nickname,password,user_type) VALUES (?,?,?,?)" //fix
	sqlParam := make([]interface{}, 0)
	default_psw := md5Gen("12345678")
	users := []string{}
	validUserType := map[string]bool{}
	validUserType["student"] = true
	if u.UserType == Admin {
		validUserType["teacher"] = true
		validUserType["admin"] = true
	}

	for i := 0; ; i++ {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.ERROR_LOG(err.Error(), c)
			c.JSON(http.StatusInternalServerError, gin.H{
				"ret": errcode.WRONG_PARAM,
				"msg": fmt.Sprintf("第%d行读取失败，请检查文件内容", i+1),
			})
			return
		}
		// skip header
		if i == 0 {
			continue
		}
		if u.UserType == Admin {
			if len(line) != 3 {
				c.JSON(http.StatusInternalServerError, gin.H{
					"ret": errcode.WRONG_PARAM,
					"msg": fmt.Sprintf("第%d行有缺失或多余字段，请检查文件内容", i+1),
				})
				return
			}
		}
		if u.UserType == Teacher {
			if len(line) != 2 {
				c.JSON(http.StatusInternalServerError, gin.H{
					"ret": errcode.WRONG_PARAM,
					"msg": fmt.Sprintf("第%d行有缺失或多余字段，请检查文件内容", i+1),
				})
				return
			}
			// 教师用户添加时，统一只能添加学生类型
			line = append(line, "student")
		}

		// 检查内容格式
		for j := 0; j < 3; j++ {
			if !utf8.ValidString(line[j]) {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"ret": errcode.WRONG_PARAM,
					"msg": fmt.Sprintf("第%d行存在非UTF-8编码字符，请将文件编码转为UTF-8", i+1),
				})
				return
			}
		}
		if !validUserType[line[2]] {
			c.JSON(http.StatusForbidden, gin.H{
				"ret": errcode.AUTH_ERR,
				"msg": fmt.Sprintf("第%d行存在不允许当前用户添加的账号类型:%s", i+1, line[2]),
			})
			return
		}
		if !useridRegexp.MatchString(line[0]) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"ret": errcode.WRONG_PARAM,
				"msg": fmt.Sprintf("第%d行用户账号非法，只允许使用32位内数字，字母，下划线", i+1),
			})
			return
		}
		users = append(users, line[0])
		// if i > 1 {
		// 	sqlTeml += ","
		// }
		// sqlTeml += "(?,?,?,?)"
		sqlParam = make([]interface{}, 0) //fix
		sqlParam = append(sqlParam, line[0], line[1], default_psw, line[2])
		tx, _ := commonlib.DB_user.Begin()
		_, err = tx.Exec(sqlTeml, sqlParam...)
		defer tx.Rollback()
		if err != nil {
			logger.ERROR_LOG(err.Error(), c)
			//fix
			mysqlErr, ok := err.(*mysql.MySQLError)
			if ok {
				logger.DEBUG_LOG(strconv.Itoa(int(mysqlErr.Number)), c)
			}
			if ok && mysqlErr.Number == 1062 {
				// Duplicate entry error, do nothing.
				// c.JSON(http.StatusInternalServerError, gin.H{
				// 	"ret": errcode.DBERR_BASE,
				// 	"msg": "账号" + line[0] + "已存在",
				// })
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"ret": errcode.DBERR_BASE,
					"msg": "插入账号时出现错误，请检查文件内容",
				})
				return
			}
		}
		tx.Commit()
	}
	// tx, _ := commonlib.DB_user.Begin()
	// _, err = tx.Exec(sqlTeml, sqlParam...)
	// defer tx.Rollback()
	// if err != nil {
	// 	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
	// 		// Duplicate entry error, do nothing.
	// 	}
	// 	logger.ERROR_LOG(err.Error(), c)
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"ret": errcode.DBERR_BASE,
	// 		"msg": "插入账号时出现错误，请检查文件内容",
	// 	})
	// 	return
	// }
	group_id, _ := strconv.Atoi(c.Query("group_id"))
	if group_id != 0 {
		valid, err := AddGroupUserDB(u, group_id, users)
		if !valid {
			c.JSON(http.StatusForbidden, gin.H{
				"ret": errcode.AUTH_ERR,
				"msg": fmt.Sprintf("导入失败，用户不允许操作组:%d", group_id),
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
	// tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "ok",
	})
}
