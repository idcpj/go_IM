package service

import (
	"errors"
	"fmt"
	"go_im/model"
	"go_im/util"
	"math/rand"
	"time"
)

type UserServer struct {
}

func (this *UserServer) Register(mobile, plainpwd, nickname, sex string) (model.User, error) {

	//手机号是否存在,存在返回已经注册
	user := model.User{
		Mobile: mobile,
	}
	_, e := DbEngin.Get(&user)
	if e != nil {
		return user, e
	}
	if user.Id > 0 {
		return user, errors.New("改手机号已经注册")
	}
	//否则插入数据
	user.Avatar = "/asset/images/avatar0.png"
	user.Sex = sex
	user.Nickname = nickname
	user.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	user.Passwd = util.MakePasswd(plainpwd, user.Salt)
	//返回新用户
	_, e = DbEngin.InsertOne(&user)
	if e != nil {
		return user, e
	}
	return user, nil
}

func (this *UserServer) Login(plainpwd, mobile string) (model.User, error) {
	user := model.User{}
	user.Mobile = mobile
	_, e := DbEngin.Get(&user)
	if e != nil {
		return user, e
	}
	if user.Id == 0 {
		return user, errors.New("user is exist")
	}

	isok := util.ValidatePasswd(plainpwd, user.Salt, user.Passwd)
	if !isok {
		return user, errors.New("password is error")
	}

	str := fmt.Sprintf("%d", time.Now().Unix())
	user.Token = util.Md5Encode(str)
	_, e = DbEngin.Id(user.Id).Cols("token").Update(user)
	if e != nil {
		return user, e
	}
	return user, nil

}

func (this *UserServer) GetUser(userId int64) (model.User,error){
	var user model.User
	user.Id= userId
	_, e := DbEngin.Get(&user)
	if e != nil {
		return user,e
	}
	return user,nil

}
