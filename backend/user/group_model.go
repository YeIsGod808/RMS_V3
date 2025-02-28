package user

import "github.com/RMS_V3/commonlib"

func GetGroupSize(group_id int) (count int) {
	commonlib.DB_user.
		QueryRow("SELECT COUNT(1) FROM t_group_user"+
			" WHERE group_id=?", group_id).
		Scan(&count)
	return
}
