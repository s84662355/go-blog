package event

import (
	"encoding/json"
 
	"GoBlog/lib/websocket"
	//"github.com/goinggo/mapstructure"
	"fmt"

)



type Event func(data map[string]interface{},session *websocket.WebsocketSession)


type EventData struct {
	Createtime int64  `json:"createtime"`
	Module     string `json:"module"`
	Event      string `json:"event"`
	Uuid       string `json:"uuid"`
	Data       map[string]interface{} `json:"data"`
}

var eventMap = make(map[string]Event)

func init(){
   
}

func Register(module string , event string,callback Event){
     eventMap[module+"_"+event] = callback
}


func Run(eventData string) bool{

	defer func(){
		//捕捉异常
		err := recover()
		if err != nil{
			fmt.Println(err)
		}
	}()

	p := &EventData{}
	err := json.Unmarshal([]byte(eventData), p)
	if err != nil{
		fmt.Println(err)
		return false
	}
	fmt.Println("事件消息体",p)

	key := p.Module + "_" + p.Event

	callback , ok := eventMap[key]

	if ok {
           session := websocket.NewWebsocketSession(p.Uuid)
           callback(p.Data,session)
           return true
	}

	fmt.Println("事件"+key+"不存在")


	return false
}