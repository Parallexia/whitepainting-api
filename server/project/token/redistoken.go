package token

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	tsgutils "github.com/typa01/go-utils"
)

var conn *redis.Pool

func init() {
	// 建立连接池
	conn = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				return nil, err
			}
			// 选择db
			//c.Do("SELECT", "token")
			return c, nil
		},
	}
}

//新建一个标记并存放在数据库
func NewToken(key string) string {
	c := conn.Get()
	defer c.Close()
	value := tsgutils.GUID()

	c.Do("set", key, value)
	return value
}

//新建一个标记并存放在数据库,设置过期时间
func NewTokenWithTTL(key string, time int) string {
	c := conn.Get()
	defer c.Close()
	value := tsgutils.GUID()

	c.Do("set", key, value)
	c.Do("EXPIRE", key, time)
	return value
}

//验证标记是否一致
func VerifyToken(key string, value string) (bool, error) {
	c := conn.Get()
	defer c.Close()

	str, err := redis.String(c.Do("GET", key))
	if err != nil {
		log.Print("Token" + err.Error())
		return false, err
	}

	if key == str {
		return true, nil
	}
	return false, nil
}

//删除标记
func DeleteToken(key string) error {
	c := conn.Get()
	defer c.Close()
	_, err := c.Do("DEL", key)
	if err != nil {
		log.Print("Token" + err.Error())
		return err
	}
	return nil
}

func ResetExpireTime(key string, time string) error {
	c := conn.Get()
	defer c.Close()

	_, err := c.Do("EXPIRE", key, time)
	if err != nil {
		log.Print("Token" + err.Error())
		return err
	}
	return nil
}

//新建一个标记，值为对应的值,默认一小时后过期
func NewTokenFromTokenValue(key string, value string) error {
	c := conn.Get()
	defer c.Close()

	c.Do("set", key, value)
	c.Do("EXPIRE", key, 60*60)
	return nil
}

//得到某个token键对应的值
func GetValueFromTokenKey(key string) (string, error) {
	c := conn.Get()
	str, err := redis.String(c.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return "", nil
		}
		log.Print("Token" + err.Error())
		return "", err
	}
	return str, nil
}
