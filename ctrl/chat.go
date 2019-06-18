package ctrl

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go_web/model"
	"go_web/service"
	"gopkg.in/fatih/set.v0"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	Id int64 `json:"id,omitempty" form:"id"` //消息ID
	//谁发的
	Userid int64 `json:"userid,omitempty" form:"userid"` //谁发的
	//什么业务
	Cmd int `json:"cmd,omitempty" form:"cmd"` //群聊还是私聊
	//发给谁
	Dstid int64 `json:"dstid,omitempty" form:"dstid"` //对端用户ID/群ID
	//怎么展示
	Media int `json:"media,omitempty" form:"media"` //消息按照什么样式展示
	//内容是什么
	Content string `json:"content,omitempty" form:"content"` //消息的内容
	//图片是什么
	Pic string `json:"pic,omitempty" form:"pic"` //预览图片
	//连接是什么
	Url string `json:"url,omitempty" form:"url"` //服务的URL
	//简单描述
	Memo string `json:"memo,omitempty" form:"memo"` //简单描述
	//其他的附加数据，语音长度/红包金额
	Amount int `json:"amount,omitempty" form:"amount"` //其他和数字相关的
}

const (
	//点对点单聊,dstid是用户ID
	CMD_SINGLE_MSG = 10
	//群聊消息,dstid是群id
	CMD_ROOM_MSG = 11
	//心跳消息,不处理
	CMD_HEART = 0
)
const (
	//文本样式
	MEDIA_TYPE_TEXT = 1
	//新闻样式,类比图文消息
	MEDIA_TYPE_News = 2
	//语音样式
	MEDIA_TYPE_VOICE = 3
	//图片样式
	MEDIA_TYPE_IMG = 4

	//红包样式
	MEDIA_TYPE_REDPACKAGR = 5
	//emoj表情样式
	MEDIA_TYPE_EMOJ = 6
	//超链接样式
	MEDIA_TYPE_LINK = 7
	//视频样式
	MEDIA_TYPE_VIDEO = 8
	//名片样式
	MEDIA_TYPE_CONCAT = 9
	//其他自己定义,前端做相应解析即可
	MEDIA_TYPE_UDEF = 100
)

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface  //把群放入其中
}

var clientMap = make(map[int64]*Node, 0)
var rw sync.RWMutex

//ws://127.0.0.1/chat>id=1&token=xxx
func Chat(w http.ResponseWriter, r *http.Request) {

	log.Println("=====================ws=====================")
	// check is Legal
	query := r.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	userId, _ := strconv.ParseInt(id, 10, 64)

	isvalida := checkToken(userId, token)
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}
	conn, e := upgrader.Upgrade(w, r, nil)
	if e != nil {
		log.Fatal(e)
		return
	}

	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//获取用户全部的群 id
	comIds,e :=contactService.SearchComunityIds(userId)
	if e != nil {
		log.Fatal(e)
		return
	}
	for _,v:=range comIds {
		node.GroupSets.Add(v)
	}

	rw.Lock()
	clientMap[userId] = node
	rw.Unlock()

	//发送逻辑
	go sendproc(node)
	//接受逻辑
	go recvproc(node)


}

func recvproc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		dispath(data)
		log.Println("data    :", string(data))
	}
}

func dispath(data []byte) {
	//解析 data 为 message

	msg := Message{}
	e := json.Unmarshal(data, &msg)
	if e != nil {
		log.Println(e)
		return
	}
	//根据 cmd 进行处理

	switch msg.Cmd {
	case CMD_SINGLE_MSG:
		sendMsg(msg.Dstid, data)
	case CMD_ROOM_MSG:
		//群聊转发逻辑
		for userId,v:=range clientMap{
			log.Println("userId=",userId)
			log.Println("msg.Id =",msg.Id)
			//过滤发送者
			if userId!=msg.Userid  && v.GroupSets.Has(msg.Dstid) {
				v.DataQueue<-data
			}
		}
	case CMD_HEART:

	}
}

func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			e := node.Conn.WriteMessage(websocket.TextMessage, data)
			if e != nil {
				log.Fatal(e)
				return
			}

		}
	}
}

func checkToken(userId int64, token string) bool {

	user := model.User{
		Id:    userId,
		Token: token,
	}
	_, e := service.DbEngin.Get(&user)
	if e != nil {
		log.Fatal(e)
		return false
	}
	if user.Mobile == "" {
		return false
	}
	return true
}

func sendMsg(userid int64, msg []byte) {
	rw.Lock()
	node, ok := clientMap[userid]
	rw.Unlock()
	if ok {
		node.DataQueue <- msg
	}
}
