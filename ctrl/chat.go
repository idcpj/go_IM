package ctrl

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go_web/model"
	"go_web/service"
	"gopkg.in/fatih/set.v0"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
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
	rw.Lock()
	clientMap[userId] = node
	rw.Unlock()

	//发送逻辑
	go sendproc(node)
	//接受逻辑
	go recvproc(node)

	sendMsg(userId, []byte("hello,world!"))

}

func recvproc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("recv     :", data)

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
