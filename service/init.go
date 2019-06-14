package service

import (
	"github.com/go-xorm/xorm"
	"go_web/model"
	"log"
)

var DbEngin *xorm.Engine

func init() {

	var e error
	driname := "mysql"
	dsName := "root:12345678@(127.0.0.1:3306)/chat?charset=utf8"
	DbEngin, e = xorm.NewEngine(driname, dsName)
	if e != nil {
		log.Fatal(e.Error())
	}
	DbEngin.ShowSQL(true)
	//设置数据库连接数
	DbEngin.SetMaxOpenConns(2)

	//创建表
	DbEngin.Sync2(
		new(model.User),
		new(model.Contact),
	)

}
