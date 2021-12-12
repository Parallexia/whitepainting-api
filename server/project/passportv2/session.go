package passportv2

import (
	"fmt"
	"log"

	tsgutils "github.com/typa01/go-utils"

	"github.com/garyburd/redigo/redis"
)

func NewUserSession(username string) (string, error) {

	c, err := DialRedis()
	if err != nil {
		return "", err
	}
	defer c.Close()

	session := tsgutils.GUID()
	expireTime := 24 * 60 * 60

	_, err = c.Do("SET", username, session, "EX", expireTime)
	if err != nil {
		log.Print("redis set failed:", err)
	}

	return session, nil
}

func VerifyUserSession(username string, session string) (bool, error) {

	c, err := DialRedis()
	if err != nil {
		return false, err
	}
	defer c.Close()

	value, err := redis.String(c.Do("GET", username))
	if err != nil {
		//如果错误是没有查询到值，错误返回nil
		if err == redis.ErrNil {
			return false, nil
		}
		log.Print("redis query failed:", err)
		return false, err
	}

	if value == session {
		return true, nil
	}
	return false, nil
}

func ResetSessionTime(username string) error {

	c, err := DialRedis()
	if err != nil {
		return err
	}
	defer c.Close()

	n, err := c.Do("EXPIRE", username, 24*60*60)
	if n == int64(1) {
		fmt.Print("success")
		return nil
	}
	return err
}

//暂时的方案
func ExpireSession(username string) error {

	c, err := DialRedis()
	if err != nil {
		return err
	}
	defer c.Close()

	n, err := c.Do("EXPIRE", username, 0.01)
	if n == int64(1) {
		fmt.Print("success")
		return nil
	}
	return err
}
