package main

import (
	_ "github.com/go-sql-driver/mysql"
	"go_web/ctrl"
	"html/template"
	"io"
	"log"
	"net/http"
)

func Refresh(w http.ResponseWriter,r *http.Request){
	registerView()
	io.WriteString(w,"Refresh success")

}
func registerView() {
	tpl, e := template.ParseGlob("view/**/*")
	if e != nil {
		log.Fatal(e.Error())
	}
	for _, v := range tpl.Templates() {
		tplname := v.Name()
		http.HandleFunc(tplname, func(w http.ResponseWriter, r *http.Request) {
			tpl.ExecuteTemplate(w, tplname, nil)
		})
	}

}

func init() {
	go registerView()
	log.SetFlags(log.Lshortfile|log.LstdFlags)
}


func main() {
	http.HandleFunc("/refresh",Refresh )
	http.HandleFunc("/user/login", ctrl.UserLogin)
	http.HandleFunc("/user/register", ctrl.UserRegister)

	http.Handle("/asset/", http.FileServer(http.Dir(".")))

	e := http.ListenAndServe(":9001", nil)
	if e != nil {
		log.Fatal("server run is error"+e.Error())
	}
	log.Panic("127.0.0.1:9001 is running")
}
