package push

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

var db *sql.DB

func InitDB() error {
	var err error
	dsn := "root" + ":" + "MysqlTest" + "@tcp(127.0.0.1:3306)/test" //记得修改

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Print(err)
		return err
	} else {
		//检测是否连接成功
		err := db.Ping()
		if err != nil {
			log.Print(err)
			return err
		}
		db.SetMaxOpenConns(10)
		fmt.Print("连接成功\n")
	}
	return nil
}

//查询项目内容
func QueryCommunity(pages int) (map[int]*CommunityInfo, error) {
	/*var id = 0
	sqlstr := `select MAX(data_id) from community;`
	result := db.QueryRow(sqlstr)
	err := result.Scan(&id)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	*/
	queryidlast := 10 * pages

	querydata := make(map[int]*CommunityInfo)
	i := 0

	sqlstr := `SELECT user_id FROM community ORDER BY data_id DESC LIMIT 10 OFFSET ?`
	results, err := db.Query(sqlstr, queryidlast)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	var data [10]CommunityInfo
	for results.Next() {

		results.Scan(&data[i].User_id)
		querydata[i] = &data[i]

		if data[i].User_id == "" {
			break
		}
		i++
	}
	defer results.Close()

	i = 0
	sqlstr = `SELECT content FROM community ORDER BY data_id DESC LIMIT 10 OFFSET ?`
	results, err = db.Query(sqlstr, queryidlast)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for results.Next() {
		results.Scan(&data[i].Content)
		querydata[i] = &data[i]
		i++
	}
	defer results.Close()

	i = 0
	sqlstr = `SELECT pic_url FROM community ORDER BY data_id DESC LIMIT 10 OFFSET ?`
	results, err = db.Query(sqlstr, queryidlast)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for results.Next() {
		results.Scan(&data[i].PicUrl)
		querydata[i] = &data[i]
		i++
	}
	defer results.Close()

	i = 0
	sqlstr = `SELECT submit_time FROM community ORDER BY data_id DESC LIMIT 10 OFFSET ?`
	results, err = db.Query(sqlstr, queryidlast)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	timehere := ""
	for results.Next() {

		results.Scan(&timehere)
		data[i].SubmitTime, err = time.Parse("2006-01-02 15:04:05", timehere)
		if err != nil {
			log.Print("communityErr")
		}
		querydata[i] = &data[i]
		i++
	}
	defer results.Close()

	return querydata, nil
}

func AddCommunityInfo(info CommunityInfo) error {
	sqlstr := `insert into community(user_id,content,pic_url,submit_time) values(?,?,?,?)`
	_, err := db.Exec(sqlstr, info.User_id, info.Content, info.PicUrl, info.SubmitTime)
	if err != nil {
		log.Print("community add item failed" + err.Error())
		return err
	}
	return nil
}
