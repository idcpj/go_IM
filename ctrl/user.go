package ctrl

import (
	"fmt"
	"go_im/args"
	"go_im/model"
	"go_im/service"
	"go_im/util"
	"log"
	"math/rand"
	"net/http"
)

var userServe = service.UserServer{}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mobile := r.PostForm.Get("mobile")
	pwd := r.PostForm.Get("passwd")
	user, e := userServe.Login(pwd, mobile)
	if e != nil {
		util.RespFail(w, e.Error())
	} else {
		util.RespOk(w, user, "")
	}

}
func UserRegister(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	mobile := r.PostForm.Get("mobile")
	if mobile == "" {
		log.Fatal("moblie not accept")
	}
	plainPwd := r.PostForm.Get("passwd")
	nickName := fmt.Sprintf("user_%06d", rand.Int31n(10000))
	sex := model.SEX_UNKNOW
	user, e := userServe.Register(mobile, plainPwd, nickName, sex)

	fmt.Printf("%+v", user)
	if e != nil {
		util.RespFail(w, e.Error())
	} else {
		util.RespOk(w, user, "")

	}
}

func UserInfo(w http.ResponseWriter, r *http.Request) {
	var user args.User
	util.Bind(r,&user)
	userInfo, e := userServe.GetUser(user.UserId)
	if e != nil {
		util.RespFail(w,e.Error())
	}else{
		util.RespOk(w,userInfo,"")
	}

}