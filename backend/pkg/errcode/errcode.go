// 错误码定义
package errcode

const (
	DBERR_BASE  string = "1000"
	DB_CONN_ERR string = "1001"
	DB_NO_ROWS  string = "1002"
	DB_DUP_ERR  string = "1062"

	SERVER_BASE    string = "2000"
	MISSING_PARAM  string = "2001"
	WRONG_PASSWORD string = "2002"
	JWT_ERR        string = "2003"
	AUTH_ERR       string = "2004"
	WRONG_PARAM    string = "2005"

	MQ_BASE     string = "3000"
	MQ_CHAN_ERR string = "3001"
)
