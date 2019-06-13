package model

import (
	"time"
)

const (
	SEX_WOMAN = "W"
	SEX_MAX = "M"
	SEX_UNKNOW = "U"
)
type User struct {
	Id       int64     `xorm:"pk autoincr bigint(64)" form:"id" json:"id"`
	Mobile   string    `xorm:"varchar(20)" form:"mobile" json:"mobile"`
	Passwd   string    `xorm:"varchar(40)" form:"passwd" json:"-"`  //pwd+salt
	Avatar   string    `xorm:"varchar(150)" form:"avatar" json:"avatar"`
	Sex      string    `xorm:"varchar(2)" form:"sex" json:"sex"`
	Nickname string    `xorm:"varchar(20)" form:"nickname" json:"nickname"`
	Salt     string    `xorm:"varchar(10)" form:"salt" json:"-"`
	Online   int       `xorm:"int(10)" form:"online" json:"online"`         //是否在线
	Token    string    `xorm:"varchar(40)" form:"token" json:"token"`			//chat?id=1&token=x
	Memo     string    `xorm:"varchar(140)" form:"memo" json:"memo"`
	Createat time.Time `xorm:"datetime created" form:"createat" json:"createat"`
}
