package ctrl

import (
	"fmt"
	"go_web/util"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	os.MkdirAll("./mnt", os.ModePerm)
}
func Upload(w http.ResponseWriter, r *http.Request) {
	UploadLocal(w, r)
}

//存储位置 ./mnt  确保已经创建好
//2. 格式  /mnt/xxx.png
func UploadLocal(w http.ResponseWriter, r *http.Request) {
	//上传的源文件
	srcfile, header, e := r.FormFile("file")
	if e != nil {
		util.RespFail(w, e.Error())
		return
	}
	suffix := ".png"
	//如果前端文件包含后缀 xxx.png
	ofilename := header.Filename
	tmp := strings.Split(ofilename, ".")
	if len(tmp) > 1 {
		suffix = "." + tmp[len(tmp)-1]
	}
	//如果前端指定了 filetype
	filetype := r.FormValue("filetype")
	if filetype != "" {
		suffix = filetype
	}
	filename := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstfile, e := os.Create("./mnt/" + filename)
	if e != nil {
		util.RespFail(w, e.Error())
		return
	}
	_, e = io.Copy(dstfile, srcfile)
	if e != nil {
		util.RespFail(w, e.Error())
		return
	}
	//将新路径拼接成 url
	url := "/mnt/" + filename
	util.RespOk(w, url, "")

}
