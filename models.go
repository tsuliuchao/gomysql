package models

import (
	"fmt"
	"log"
	"strings"
	"math/rand"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/gin-blog/pkg/setting"
	"github.com/go-ini/ini"
	"time"
)

var db *gorm.DB

// GetConn("base", true, 2018)
// @database_config db数据库配置名称 
// @is_read 是否只读链接
func GetConn(database_config string,is_read bool) () {
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)
	sources, err := setting.Cfg.GetSection(database_config)
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}
    hosts := sources.Key("HOST").String()
	hosts_array := strings.Split(hosts, ",")
	fmt.Println(hosts_array)
	if len(hosts_array) == 0 {
		//return nil, errors.New("empty servers list")
		log.Fatal(2, "empty servers list")

	}
	if is_read==false {
		host = hosts_array[0]
	}else{
		max:=len(hosts_array)
		host = hosts_array[RandInt64(1,int64(max))]
	}
	dbType = sources.Key("TYPE").String()
	dbName = sources.Key("NAME").String()
	user = sources.Key("USER").String()
	password = sources.Key("PASSWORD").String()
	tablePrefix = sources.Key("TABLE_PREFIX").String()
	db, err = gorm.Open(dbType,fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",user,password,host,dbName))

	if err != nil {
		log.Println(err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}
	db.SingularTable(true)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func RandInt64(min, max int64) int64 {
    rand.Seed(time.Now().UnixNano())
    return min + rand.Int63n(max-min)
}

func CloseDB() {
	defer db.Close()
}



















