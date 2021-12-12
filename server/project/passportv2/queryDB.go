package passportv2

import (
	"log"
)

//从数据库查询电子邮件
func QueryEmail(username string) string {
	var email string
	sqlstr := `select email from users where username=?`
	result := db.QueryRow(sqlstr, username)
	err := result.Scan(&email)
	if err != nil {
		log.Print(err)
		return ""
	}
	return email
}

//从数据库查询用户名
func QueryUsername(email string) string {
	var username string
	sqlstr := `select username from users where email=?`
	result := db.QueryRow(sqlstr, email)
	err := result.Scan(&username)
	if err != nil {
		log.Print(err)
		return ""
	}
	return username
}

//有可能有逻辑问题
func QueryUsernameIsExist(username string) bool {
	var usernameSearch string
	sqlstr := `select username from users where username=?`
	result := db.QueryRow(sqlstr, username)
	err := result.Scan(&usernameSearch)
	if usernameSearch == "" {
		return false
	}

	if err != nil {
		log.Print(err)
		return true
	}

	return true
}

//有可能有逻辑问题
func QueryEmailIsExist(email string) bool {
	var usernameSearch string
	sqlstr := `select email from users where email=?`
	result := db.QueryRow(sqlstr, email)
	err := result.Scan(&usernameSearch)
	if usernameSearch == "" {
		return false
	}

	if err != nil {
		log.Print(err)
		return true
	}

	return true
}

//从数据库查询密码
func QueryPasswd(username string) []byte {
	var passwd []byte
	sqlstr := `select password from users where username=?`
	result := db.QueryRow(sqlstr, username)
	err := result.Scan(&passwd)
	if err != nil {
		log.Print(err)
		return nil
	}
	return []byte(passwd)
}

//通过邮箱查询用户id
func QueryIdByEmail(email string) string {
	var id string
	sqlstr := `select data_id from users where email=?`
	result := db.QueryRow(sqlstr, email)
	err := result.Scan(&id)
	if err != nil {
		log.Print(err)
		return ""
	}
	return id
}

//通过用户名查询用户id
func QueryIdByUsername(username string) string {
	var id string
	sqlstr := `select data_id from users where username=?`
	result := db.QueryRow(sqlstr, username)
	err := result.Scan(&id)
	if err != nil {
		if err == nil {
			return ""
		}
		log.Print(err)
		return ""
	}
	return id
}

//从数据库查询盐
func QuerySalt(username string) []byte {
	var salt []byte
	sqlstr := `select salt from users where username=?`
	result := db.QueryRow(sqlstr, username)
	err := result.Scan(&salt)
	if err != nil {
		log.Print(err)
		return nil
	}
	return salt
}

//通过id查询用户名
func QueryUsernameById(id string) string {
	var username string
	sqlstr := `select username from users where data_id=?`
	result := db.QueryRow(sqlstr, id)
	err := result.Scan(&username)
	if err != nil {
		if err == nil {
			return ""
		}
		log.Print(err)
		return ""
	}
	return username
}

func QueryLastID() (string, error) {
	id := ""
	sqlstr := `select MAX(data_id) from users;`
	result := db.QueryRow(sqlstr)
	err := result.Scan(&id)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return id, nil
}
