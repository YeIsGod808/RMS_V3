package user

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/RMS_V3/logger"

	"github.com/RMS_V3/commonlib"
	"github.com/RMS_V3/pkg/errcode"

	"github.com/gin-gonic/gin"
)

func CreateGroup(c *gin.Context) {
	u, ok := CheckUserPermission(c.Query("token"), Teacher, c)
	if !ok {
		return
	}
	res, err := commonlib.DB_user.Exec("INSERT INTO t_group SET owner=?,name=?", u.Id, c.Query("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.DBERR_BASE,
			"msg": "fail",
		})
		return
	}
	id, _ := res.LastInsertId()
	c.JSON(http.StatusOK, gin.H{
		"ret":      "0",
		"msg":      "succ",
		"group_id": id,
	})
}

func AddGroupUser(c *gin.Context) {
	u, ok := CheckUserPermission(c.Query("token"), Teacher, c)
	if !ok {
		return
	}
	group_id, _ := strconv.Atoi(c.Query("group_id"))
	var postJson struct {
		Users []string
	}
	c.BindJSON(&postJson)
	valid, err := AddGroupUserDB(u, group_id, postJson.Users)
	if !valid {
		c.JSON(http.StatusForbidden, gin.H{
			"ret": errcode.AUTH_ERR,
			"msg": "not owner of group",
		})
		return
	}
	if err != nil {
		if commonlib.IsDbDupErr(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"ret": errcode.DB_DUP_ERR,
				"msg": "错误：试图导入组内已有账号",
			})
			return
		}
		logger.ERROR_LOG(err.Error(), c)
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.DBERR_BASE,
			"msg": "fail",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": "0",
		"msg": "succ",
	})
}

func AddGroupUserDB(u *User, group_id int, users []string) (valid bool, err error) {
	var tmp int
	err = commonlib.DB_user.QueryRow("SELECT 1 FROM t_group"+
		" WHERE group_id=? AND owner=?", group_id, u.Id).Scan(&tmp)
	if err != nil {
		valid = false
		return
	}
	// parse request: end
	valid = true
	sqlTmpl := "INSERT INTO t_group_user VALUES "
	sqlParam := make([]interface{}, 0, len(users)*2)
	for i, user_id := range users {
		if i != 0 {
			sqlTmpl += ","
		}
		sqlTmpl += "(?,?)"
		sqlParam = append(sqlParam, group_id, user_id)
	}
	_, err = commonlib.DB_user.Exec(sqlTmpl, sqlParam...)
	return
}

func deleteGroupUsersDB(group_id int, users []string) (err error) {
	sqlTmpl := "DELETE FROM t_group_user WHERE (group_id,user_id) in ("
	sqlParam := make([]interface{}, 0, len(users)*2)
	for i, user_id := range users {
		if i != 0 {
			sqlTmpl += ","
		}
		sqlTmpl += "(?,?)"
		sqlParam = append(sqlParam, group_id, user_id)
	}
	sqlTmpl += ")"
	_, err = commonlib.DB_user.Exec(sqlTmpl, sqlParam...)
	return
}

func DeleteGroupUser(c *gin.Context) {
	u, ok := CheckUserPermission(c.Query("token"), Teacher, c)
	if !ok {
		return
	}
	group_id, _ := strconv.Atoi(c.Query("group_id"))
	var postJson struct {
		Users []string
	}
	c.BindJSON(&postJson)
	var tmp int
	err := commonlib.DB_user.QueryRow("SELECT 1 FROM t_group"+
		" WHERE group_id=? AND owner=?", group_id, u.Id).Scan(&tmp)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"ret": errcode.AUTH_ERR,
			"msg": "not owner of group",
		})
		return
	}
	// parse request: end
	err = deleteGroupUsersDB(group_id, postJson.Users)
	if err != nil {
		logger.ERROR_LOG(err.Error(), c)
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.DBERR_BASE,
			"msg": "fail",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": "0",
		"msg": "succ",
	})
}

type UserIdNick struct {
	User_id  string `json:"user_id"`
	Nickname string `json:"nickname"`
}
type GroupInfo struct {
	Owner    UserIdNick   `json:"owner"`
	Name     string       `json:"name"`
	Group_id int          `json:"group_id"`
	Users    []UserIdNick `json:"users"`
}

func GetGroupDB(group_id int) (group_info GroupInfo, err error) {
	group_info.Group_id = group_id
	err = commonlib.DB_user.
		QueryRow("SELECT owner,name FROM t_group"+
			" WHERE group_id=?", group_id).
		Scan(&group_info.Owner.User_id, &group_info.Name)
	if err != nil {
		return
	}
	err = commonlib.DB_user.
		QueryRow("SELECT nickname FROM t_account"+
			" WHERE user_id=?", group_info.Owner.User_id).
		Scan(&group_info.Owner.Nickname)
	if err != nil {
		return
	}
	rows, err := commonlib.DB_user.Query("SELECT user_id,nickname FROM t_group_user"+
		" left join (SELECT user_id,nickname FROM t_account)b using(user_id)"+
		" WHERE group_id=?", group_id)
	if err != nil {
		return
	}
	for rows.Next() {
		// nickname可能会有null，所以用byte解决
		var user_id, nickname sql.RawBytes
		rows.Scan(&user_id, &nickname)
		group_info.Users = append(group_info.Users,
			UserIdNick{string(user_id), string(nickname)})
	}
	err = rows.Close()
	return
}

// only the owner or member of group can get user list
func GetGroupUser(c *gin.Context) {
	u, ok := CheckUserPermission(c.Query("token"), Student, c)
	if !ok {
		return
	}
	group_id, _ := strconv.Atoi(c.Query("group_id"))
	var tmp int
	err := commonlib.DB_user.QueryRow("SELECT 1 FROM t_group_user"+
		" WHERE group_id=? AND user_id=?", group_id, u.Id).Scan(&tmp)
	if err != nil {
		err = commonlib.DB_user.QueryRow("SELECT 1 FROM t_group"+
			" WHERE group_id=? AND owner=?", group_id, u.Id).Scan(&tmp)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"ret": errcode.AUTH_ERR,
				"msg": "not owner or member of group",
			})
			return
		}
	}
	// parse request end

	group_info, err := GetGroupDB(group_id)
	if err != nil {
		logger.ERROR_LOG(err.Error(), c)
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.DBERR_BASE,
			"msg": "fail",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret":   "0",
		"group": group_info,
	})
}

func DeleteGroup(c *gin.Context) {
	u, ok := CheckUserPermission(c.Query("token"), Teacher, c)
	if !ok {
		return
	}
	group_id, _ := strconv.Atoi(c.Query("group_id"))
	var tmp int
	err := commonlib.DB_user.QueryRow("SELECT 1 FROM t_group"+
		" WHERE group_id=? AND owner=?", group_id, u.Id).Scan(&tmp)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"ret": errcode.AUTH_ERR,
			"msg": "not owner of group",
		})
		return
	}
	// parse end
	_, err = commonlib.DB_user.Exec("DELETE FROM t_group_user WHERE group_id=?",
		group_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.DBERR_BASE,
			"msg": "fail when delete users",
		})
		return
	}
	_, err = commonlib.DB_user.Exec("DELETE FROM t_group WHERE group_id=?",
		group_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.DBERR_BASE,
			"msg": "fail when delete group",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": "0",
		"msg": "succ",
	})
}

func GetGroupList(c *gin.Context) {
	u, ok := CheckUserPermission(c.Query("token"), Student, c)
	if !ok {
		return
	}
	type group_info struct {
		Group_id int    `json:"group_id"`
		Name     string `json:"name"`
	}
	var groups struct {
		Owned  []group_info `json:"owned"`
		Joined []group_info `json:"joined"`
	}
	tmpGroup := group_info{}
	rows, err := commonlib.DB_user.Query("SELECT group_id,name"+
		" FROM t_group WHERE owner=?", u.Id)
	if err != nil {
		logger.ERROR_LOG(err.Error(), c)
	} else {
		for rows.Next() {
			rows.Scan(&tmpGroup.Group_id, &tmpGroup.Name)
			groups.Owned = append(groups.Owned, tmpGroup)
		}
		rows.Close()
	}
	/* user_id 对应的 group_id, name */
	rows, err = commonlib.DB_user.Query("SELECT group_id,name"+
		" FROM t_group_user inner join"+
		"(SELECT group_id,name FROM t_group) b using(group_id)"+
		" WHERE user_id=?", u.Id)
	if err != nil {
		logger.ERROR_LOG(err.Error(), c)
	} else {
		for rows.Next() {
			rows.Scan(&tmpGroup.Group_id, &tmpGroup.Name)
			groups.Joined = append(groups.Joined, tmpGroup)
			/* 只能查到自己所在的group */
		}
		rows.Close()
	}
	c.JSON(http.StatusOK, gin.H{
		"ret":    "0",
		"groups": groups,
	})
}

func GetAllGroup(c *gin.Context) {
	_, ok := CheckUserPermission(c.Query("token"), Admin, c)
	if !ok {
		return
	}
	type group_info struct {
		Group_id int    `json:"group_id"`
		Name     string `json:"name"`
	}
	res := make([]group_info, 0)
	tmpGroup := group_info{}
	rows, err := commonlib.DB_user.Query("SELECT group_id,name FROM t_group")
	if err != nil {
		logger.ERROR_LOG(err.Error(), c)
	} else {
		for rows.Next() {
			rows.Scan(&tmpGroup.Group_id, &tmpGroup.Name)
			res = append(res, tmpGroup)
		}
		rows.Close()
	}
	c.JSON(http.StatusOK, gin.H{
		"ret":    "0",
		"groups": res,
	})
}

func EditGroupName(c *gin.Context) {
	u, ok := CheckUserPermission(c.Query("token"), Teacher, c)
	if !ok {
		return
	}
	group_id, _ := strconv.Atoi(c.Query("group_id"))
	name := c.Query("name")
	var tmp int
	err := commonlib.DB_user.QueryRow("SELECT 1 FROM t_group"+
		" WHERE group_id=? AND owner=?", group_id, u.Id).Scan(&tmp)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"ret": errcode.AUTH_ERR,
			"msg": "not owner of group",
		})
		return
	}
	// parse end
	_, err = commonlib.DB_user.Exec("UPDATE t_group SET name=?"+
		" WHERE group_id=?", name, group_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ret": errcode.DBERR_BASE,
			"msg": "fail",
		})
	}
}

func IsGroupOwner(u *User, group_id int, c *gin.Context) (ok bool) {
	var tmp int
	if commonlib.DB_user.
		QueryRow("SELECT 1 FROM t_group"+
			" WHERE group_id=? AND owner=?", group_id, u.Id).
		Scan(&tmp) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"ret": errcode.AUTH_ERR,
			"msg": "not owner of group",
		})
		return false
	}
	return true
}

// check if user is member of group,
// return false and response 403 if not
func IsGroupOwnerOrMember(u *User, group_id int, c *gin.Context) (ok bool) {
	var tmp int
	err := commonlib.DB_user.QueryRow("SELECT 1 FROM t_group_user"+
		" WHERE group_id=? AND user_id=?", group_id, u.Id).Scan(&tmp)
	if err != nil {
		err = commonlib.DB_user.QueryRow("SELECT 1 FROM t_group"+
			" WHERE group_id=? AND owner=?", group_id, u.Id).Scan(&tmp)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"ret": errcode.AUTH_ERR,
				"msg": "not owner or member of group",
			})
			return false
		}
	}
	return true
}
