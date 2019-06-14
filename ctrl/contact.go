package ctrl

import (
	"go_web/args"
	"go_web/service"
	"go_web/util"
	"net/http"
)

var contactService service.ContactService

func Addfriend(w http.ResponseWriter, r *http.Request) {

	var arg args.ContactArg
	//对象绑定
	util.Bind(r, &arg)

	err := contactService.AddFriend(arg.Userid, arg.Dstid)

	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "好友添加成功")
	}
}

func Loadfriend(w http.ResponseWriter, r *http.Request) {

	var arg args.ContactArg
	//对象绑定
	util.Bind(r, &arg)

	users, e := contactService.SearchFriend(arg.Userid)
	if e != nil {
		util.RespFail(w, "no user lists")
	}

	util.RespOkList(w, users, int64(len(users)))

}
