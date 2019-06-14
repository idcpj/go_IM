package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type H struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Rows  interface{} `json:"rows,omitempty"`
	Total interface{} `json:"total,omitempty"`
	Data  interface{} `json:"data,omitempty"` //omitempty 转为 json 为空就不带上 data
}
type Hs struct {
	Code int         `json:"code"`
	Rows string      `json:"msg"`
	Data interface{} `json:"data,omitempty"` //omitempty 转为 json 为空就不带上 data
}

func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}
func RespOk(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0, data, msg)

}
func RespOkList(w http.ResponseWriter, lists interface{}, total int64) {
	RespList(w, 0, lists, total)

}
func RespList(w http.ResponseWriter, code int, data interface{}, total int64) {
	w.Header().Set("Content-Type", "application/json")
	//设置200状态
	w.WriteHeader(http.StatusOK)
	//测试 100
	//20
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	//将结构体转化成JSOn字符串
	ret, err := json.Marshal(h)
	if err != nil {
		log.Println(err.Error())
	}
	//输出
	w.Write(ret)
}
func Resp(w http.ResponseWriter, code int, data interface{}, err string) {
	w.Header().Set("Content-Type", "application/json")
	h := H{
		Msg:  err,
		Code: code,
		Data: data,
	}
	bytes, _ := json.Marshal(h)
	w.Write(bytes)
}
