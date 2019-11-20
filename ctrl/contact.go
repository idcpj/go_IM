package ctrl

import (
	"go_im/args"
	"go_im/service"
	"go_im/util"
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


func LoadCommunity(w http.ResponseWriter, req *http.Request){
	var arg args.ContactArg
	//如果这个用的上,那么可以直接
	util.Bind(req,&arg)
	comunitys := contactService.SearchComunity(arg.Userid)
	util.RespOkList(w,comunitys,int64(len(comunitys)))
}
func JoinCommunity(w http.ResponseWriter, req *http.Request){
	var arg args.ContactArg

	//如果这个用的上,那么可以直接
	util.Bind(req,&arg)
	err := contactService.JoinCommunity(arg.Userid,arg.Dstid)

	if err!=nil{
		util.RespFail(w,err.Error())
	}else {
		addGroupId(arg.Userid,arg.Dstid)
		util.RespOk(w,nil,"")
	}
}

func AddCommunity(w http.ResponseWriter, req *http.Request){
	var arg args.ContactArg
	util.Bind(req,&arg)
	group,e := contactService.CreateGroup(arg.Userid, arg.GroupName)
	if e != nil {
		util.RespFail(w,e.Error())
	}else{
		util.RespOk(w,group,"")
	}
}

func addGroupId(userId int64, groupId int64) {
	rw.Lock()
	node,ok:=clientMap[userId]
	if ok {
		node.GroupSets.Add(groupId)
	}
	rw.Unlock()
}