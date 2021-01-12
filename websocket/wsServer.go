package websocket

import (
	"net/http"
	"fmt"
	"strconv"
	"time"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
	"sync"
	"GoBlog/lib/helper"
)

type wsServer struct {
   server *http.Server
   port uint16
   ip string
   connMap map[string]net.Conn
   connlock sync.RWMutex
   connection uint32
   pathMap map[string]func (conn net.Conn)
}

func (l *wsServer) Start(ip string,port uint16){
    l.port = port
    l.ip = ip
    l.connMap = make(map[string]net.Conn)
}

func (l *wsServer) initServer() {
	l.server = new(http.Server)
	l.server.Addr = l.ip+":"+ strconv.Itoa(int(l.port))
	l.server.ReadTimeout = 10 * time.Second
	l.server.WriteTimeout = 10 * time.Second
	l.server.MaxHeaderBytes = 1 << 20
    l.server.Handler = l
}

func (l *wsServer)  addConn(conn net.Conn) (res  bool){
       l.connlock.Lock()
       defer func(){
		   l.connlock.Unlock()
		   if r := recover(); r != nil {
			   fmt.Printf("捕获到的错误：%s\n", r)
			   res = false
		   }
	   }()
	   res = true
       uuid := helper.GetRandomString(10) + helper.GetLocalIP()
       l.connMap[uuid] = conn
       return res

}




func (l *wsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	switch r.URL.Path {
	case "/":
		fmt.Println("hello")
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not found: %s\n", r.URL)
	}
}