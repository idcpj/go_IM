package service

import (
	"errors"
	"go_web/model"
)

type ContactService struct {
}

func (this *ContactService) AddFriend(userid, dstid int64) error {

	contact := model.Contact{}
	if userid == dstid {
		return errors.New("not add self for friend")
	}
	//是否已经添加自己为好友
	_, e := DbEngin.Where("ownerid=?", userid).And("dstobj=?", dstid).And("cate=?", model.CONCAT_CATE_USER).Get(&contact)
	if e != nil {
		return e
	}
	if contact.Id > 0 {
		return errors.New("user account is add the friends")
	}
	session := DbEngin.NewSession()
	session.Begin()
	//插入自己的
	_, e1 := session.InsertOne(model.Contact{
		Ownerid: userid,
		Dstobj:  dstid,
		Cate:    model.CONCAT_CATE_USER,
	})
	//插入对方的

	_, e2 := session.InsertOne(model.Contact{
		Ownerid: dstid,
		Dstobj:  userid,
		Cate:    model.CONCAT_CATE_USER,
	})
	if e1 == nil && e2 == nil {
		session.Commit()
		return nil
	} else {
		session.Rollback()
		if e1 != nil {
			return e1
		} else {
			return e2
		}
	}

}
func (this *ContactService) SearchFriend(userid int64) ([]model.User, error) {
	//var
	userids := make([]string, 0)
	e := DbEngin.Table(new(model.Contact)).Where("ownerid=?", userid).And("cate=?", model.CONCAT_CATE_USER).Cols("dstobj").Find(&userids)
	if e != nil {
		return nil, e
	}
	users := make([]model.User, 0)
	e = DbEngin.In("id", userids).Find(&users)
	if e != nil {
		return nil, e
	}
	return users, nil
}

func (this *ContactService) SearchComunityIds(userId int64) ([]int64, error) {
	conids := make([]int64, 0)
	e := DbEngin.Table(new(model.Contact)).Where("ownerid =? and cate=?", userId, model.CONCAT_CATE_COMUNITY).Cols("dstobj").Find(&conids)
	return conids, e
}

func (service *ContactService) SearchComunity(userId int64) ([]model.Community) {
	comIds := make([]int64, 0)

	DbEngin.Table(new(model.Contact)).Where("ownerid = ? and cate = ?", userId, model.CONCAT_CATE_COMUNITY).Cols("dstobj").Find(&comIds)

	if len(comIds) == 0 {
		return nil
	}
	coms := make([]model.Community, 0)
	DbEngin.In("id", comIds).Find(&coms)
	return coms
}

//加群
func (service *ContactService) JoinCommunity(userId, comId int64) error {
	cot := model.Contact{
		Ownerid: userId,
		Dstobj:  comId,
		Cate:    model.CONCAT_CATE_COMUNITY,
	}
	DbEngin.Get(&cot)
	if (cot.Id == 0) {
		_, err := DbEngin.InsertOne(cot)
		return err
	} else {
		return nil
	}
}

func (this *ContactService) CreateGroup(userId int64, groupName string) (model.Community, error) {
	group := model.Community{}
	group.Ownerid = userId
	group.Name = groupName
	group.Cate = 2
	group.Memo = ""
	group.Icon = "/asset/images/avatar0.png"
	_, e := DbEngin.InsertOne(&group)
	return group, e
}
