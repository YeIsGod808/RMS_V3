package commonlib

func IsDbDupErr(err error) bool {
	if len(err.Error()) >= 10 && err.Error()[0:10] == "Error 1062" {
		return true
	}
	return false
}
