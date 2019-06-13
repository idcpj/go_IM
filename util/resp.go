package util

import (
	"encoding/json"
	"net/http"
)


type H struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"` //omitempty 转为 json 为空就不带上 data
}
func RespFail(w http.ResponseWriter,msg string){
	Resp(w,-1,nil,msg)
}
func RespOk(w http.ResponseWriter,data interface{},msg string){
	Resp(w,0,data,msg)

}
func Resp(w http.ResponseWriter,code int, data interface{}, err string) {
	w.Header().Set("Content-Type","application/json")
	h := H{
		Msg:  err,
		Code: code,
		Data: data,
	}
	bytes, _ := json.Marshal(h)
	w.Write(bytes)
}

