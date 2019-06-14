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
